#![allow(dead_code)]
#![allow(unused_variables)]

use rand::prelude::*;
use rand::rngs::mock::StepRng;
use shuffle::irs::Irs;
use shuffle::shuffler::Shuffler;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fs;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::iter::FromIterator;

#[derive(Copy, Clone)]

struct Key {
    //where b and c are the 'keyed' grids, a and d are filled A-Z
    a: [u8; 25],
    b: [u8; 25],
    c: [u8; 25],
    d: [u8; 25],
}

#[derive(Debug)]
struct Coords {
    row: i32,
    column: i32,
}

struct Solution {
    plaintext: String,
    key: Key,
}

const LETTERS: [u8; 25] = [
    65, 66, 67, 68, 69, 70, 71, 72, 73, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89,
    90,
]; //excludes J

fn format_text(original: Vec<u8>) -> Vec<u8> {
    //format to remove spaces and characters that are not capital letters
    let mut result: Vec<u8> = Vec::new();
    for i in 0..original.len() {
        if original[i] >= 65 && original[i] <= 90 {
            result.push(original[i])
        }
    }
    result
}

fn scan_tetragrams() -> HashMap<Vec<u8>, i32> {
    //scans file of tetragram frequencies and loads into hashmap
    let mut tg: HashMap<Vec<u8>, i32> = HashMap::new();
    let file = File::open("../tetragrams.txt").expect("Failed to load tetragram file");
    let reader = BufReader::new(&file);
    for line in reader.lines() {
        let ln: &str = &line.unwrap_or(String::from("Invalid line"));
        let fields: Vec<&str> = ln.split(", ").collect();
        let n: Result<i32, _> = fields[1].parse();
        let num: i32 = match n {
            Ok(n) => n,
            Err(_) => panic!("Failed to convert tetragram frequency to i32"),
        };
        tg.insert(fields[0].as_bytes().to_vec(), num);
    }
    tg
}

fn gen_reference(n: usize) -> Vec<u8> {
    //generate reference text of same length as ciphertext to compare tetregram scores
    let r: Vec<u8> = format_text(read("../referenceText.txt"));
    let mut ref_text: Vec<u8> = Vec::new();
    for i in 0..n {
        ref_text.push(r[i]);
    }
    ref_text
}

fn scan_dict() -> Vec<Vec<u8>> {
    //scans dictionary file into a vector
    //very inefficient with conversions, but this is only called once
    let words: Vec<String> = String::from_utf8(read("../challenge_word_list.txt").clone())
        .unwrap()
        .split(" ")
        .collect::<Vec<&str>>()
        .iter()
        .map(|&s| s.into())
        .collect(); //this was painful lmao
    let mut dict: Vec<Vec<u8>> = Vec::new();
    for mut w in words {
        w.retain(|c| {
            let mut seen: HashSet<char> = HashSet::new();
            let s = seen.contains(&c);
            seen.insert(c);
            !s
        });
        if w.contains("J") || w.len() == 1 {
            continue;
        }
        dict.push(w.into_bytes())
    }
    dict[..dict.len() - 1].to_vec() //last value is a newline character
}

fn randomise() -> [u8; 25] {
    //randomises grid with letters A-Z excluding J
    let mut g: Vec<u8> = LETTERS.to_vec();
    let mut rng = StepRng::new(0, 1);
    let mut irs = Irs::default();
    let _ = irs.shuffle(&mut g, &mut rng);
    let t: Result<[u8; 25], _> = g.try_into();
    match t {
        Ok(arr) => arr,
        Err(_) => panic!("Failed to randomise grid"),
    }
}

fn read(path: &str) -> Vec<u8> {
    let f: Result<Vec<u8>, std::io::Error> = fs::read(path);
    match f {
        Ok(s) => s,
        Err(_) => panic!("Couldn't find file"),
    }
}

fn score(text: &Vec<u8>, tetragrams: &HashMap<Vec<u8>, i32>) -> i32 {
    //fitness score based on tetragram frequency
    let mut total: i32 = 0;
    for i in 0..text.len() - 3 {
        let v: Vec<u8> = text[i..i + 4].to_vec();
        if tetragrams.contains_key(&v) {
            total += tetragrams[&v];
        }
    }
    total
}

fn fill_grid(keyword: Vec<u8>) -> [u8; 25] {
    //fills grid based on keyword
    let mut seen: HashSet<u8> = HashSet::from_iter(keyword.clone()); //collect letters from keyword
    let last_idx = LETTERS
        .iter()
        .position(|&x| x == *keyword.last().unwrap())
        .unwrap(); //continues filling up grid with next letter after last letter in keyword
    let mut grid: Vec<u8> = keyword; //fill with letters of keyword
    for i in last_idx..last_idx + 25 {
        let idx = i % 25;
        if !seen.contains(&LETTERS[idx]) {
            seen.insert(LETTERS[idx]);
            grid.push(LETTERS[idx]);
        }
    }
    let t: Result<[u8; 25], _> = grid.try_into();
    match t {
        Ok(arr) => arr,
        Err(_) => panic!("Failed to fill grid"),
    }
}

fn coordinates(grid_index: i32) -> Coords {
    //coords in grid from index
    let row: i32 = grid_index / 5;
    let column: i32 = grid_index - row * 5;
    Coords { row, column }
}

fn grid_index(row: i32, column: i32) -> i32 {
    //grid index from coords
    row * 5 + column
}

fn decipher(text: &Vec<u8>, k: Key) -> Option<Vec<u8>> {
    //decipher text given key
    let mut deciphered: Vec<u8> = Vec::new();
    for i in (0..text.len() - 1).step_by(2) {
        let v: Vec<u8> = (text[i..i + 2]).to_vec();
        let i1: usize = k.b.iter().position(|&x| x == v[0])?;
        let i2: usize = k.c.iter().position(|&x| x == v[1])?;
        //indices of letters in keyed grids

        let coords1: Coords = coordinates(i1.try_into().ok()?);
        let coords2: Coords = coordinates(i2.try_into().ok()?);
        //coords in key grids

        let d_i1: usize = grid_index(coords1.row, coords2.column)
            .try_into()
            .expect("Failed to get grid index 1"); //coords of deciphered letters
        let d_i2: usize = grid_index(coords2.row, coords1.column)
            .try_into()
            .expect("Failed to get grid index 2");
        //indixes in a and d

        deciphered.push(k.a[d_i1]);
        deciphered.push(k.d[d_i2]);
    }
    Some(deciphered)
}

fn dictionary_attack(ciphertext: &Vec<u8>) {
    let dict = scan_dict();
    let tetragrams = scan_tetragrams();
    let ref_score = score(&gen_reference(ciphertext.clone().len()), &tetragrams);
    for word in &dict {
        for word2 in &dict {
            let k: Key = Key {
                a: LETTERS,
                b: fill_grid(word.to_vec()),
                c: fill_grid(word2.to_vec()),
                d: LETTERS,
            };
            let p = decipher(ciphertext, k).unwrap();
            if score(&p, &tetragrams) * 10 >= ref_score * 8 {
                println!("{}", String::from_utf8(p).unwrap());
            }
        }
    }
}

fn mutate(k: &mut Key) {
    //randomly swaps two elements in one of the two grids
    let mut rng = rand::thread_rng();
    let x: f64 = rng.gen(); //rng choosing which grid to mutate
    if x > 0.5 {
        let i1 = rng.gen_range(0..25);
        let mut i2 = rng.gen_range(0..25);
        while i1 == i2 {
            i2 = rng.gen_range(0..25);
        }
        let temp = k.b[i1];
        k.b[i1] = k.b[i2];
        k.b[i2] = temp
    } else {
        let i1 = rng.gen_range(0..25);
        let mut i2 = rng.gen_range(0..25);
        while i1 == i2 {
            i2 = rng.gen_range(0..25);
        }
        let temp = k.c[i1];
        k.c[i1] = k.c[i2];
        k.c[i2] = temp
    }
}

fn acceptance_probability(delta_e: f64, temp: f64) -> f64 {
    //crude custom acceptance probability function
    //based on idea that high delta_e should decrease p, high temp should increase p
    (1.0 / delta_e) * temp
}

fn simulated_annealing(ciphertext: &Vec<u8>, max_constant: i32, max_temp: f64, k: f64) -> Solution {
    //simulated annealing to aim to avoid local minima
    let mut iterations_count = 0;
    let mut temp = max_temp;
    let tetragrams: HashMap<Vec<u8>, i32> = scan_tetragrams();
    let mut new_key = Key {
        a: LETTERS,
        b: randomise(),
        c: randomise(),
        d: LETTERS,
    };
    let mut constant: i32 = 0;
    while constant < max_constant {
        iterations_count += 1;
        let old_key = new_key.clone();
        mutate(&mut new_key);
        let s_old: f64 = 1.0;
        let s_new: f64 = score(&decipher(ciphertext, new_key).unwrap(), &tetragrams) as f64
            / score(&decipher(ciphertext, old_key).unwrap(), &tetragrams) as f64;
        let delta_e: f64 = s_old / s_new;
        if delta_e > 1.0 {
            //old key is better
            let mut rng = rand::thread_rng();
            let x: f64 = rng.gen();
            let p = acceptance_probability(delta_e, temp);
            if x < p {
                constant = 0;
            } else {
                new_key = old_key;
                constant += 1;
            }
        } else {
            constant = 0;
        }
        temp *= k;
    }
    Solution {
        plaintext: String::from_utf8(decipher(ciphertext, new_key).unwrap()).unwrap(),
        key: new_key,
    }
}

fn renneal(ciphertext: &Vec<u8>, max_constant: i32, max_temp: f64, k: f64) -> String {
    //simulated annealing until the same solution is reached twice - runtime pretty long but should
    //be basically always correct
    let mut solutions: HashSet<String> = HashSet::new();
    loop {
        let res: Solution = simulated_annealing(ciphertext, max_constant, max_temp, k);
        let p: String = res.plaintext;
        println!("{}", p);
        if solutions.contains(&p) {
            //repeated solution is very likely to be correct
            println!("Key 1: {}", String::from_utf8(res.key.b.to_vec()).unwrap());
            println!("Key 2: {}", String::from_utf8(res.key.c.to_vec()).unwrap());
            return p;
        }
        solutions.insert(p);
        println!("{}", solutions.len());
    }
}

fn reference_annealing(ciphertext: &Vec<u8>, max_constant: i32, max_temp: f64, k: f64) -> String {
    //simulated annealing until a solution with score > 80% of reference score is found
    let tetragrams = scan_tetragrams();
    let ref_score = score(&gen_reference(ciphertext.len()), &tetragrams);
    loop {
        let res: Solution = simulated_annealing(ciphertext, max_constant, max_temp, k);
        let p: String = res.plaintext;
        if score(&p.as_bytes().to_vec(), &tetragrams) * 10 >= ref_score * 8 {
            println!("Key 1: {}", String::from_utf8(res.key.b.to_vec()).unwrap());
            println!("Key 2: {}", String::from_utf8(res.key.c.to_vec()).unwrap());
            return p;
        }
    }
}

fn main() {
    let original = read("src/ciphertext.txt");
    let ciphertext: Vec<u8> = format_text(original);
    println!("{}", reference_annealing(&ciphertext, 1000, 1.0, 0.95));
}
