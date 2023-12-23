use rand::prelude::*;
use rand::rngs::mock::StepRng;
use shuffle::irs::Irs;
use shuffle::shuffler::Shuffler;
use std::collections::HashMap;
use std::fs;
use std::fs::File;
use std::io::{BufRead, BufReader};

const LETTERS: [u8; 26] = [
    65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
    89, 90,
];

fn format_text(original: Vec<u8>) -> Vec<u8> {
    let mut result: Vec<u8> = Vec::new();
    for i in 0..original.len() {
        if original[i] >= 65 && original[i] <= 90 {
            result.push(original[i])
        }
    }
    result
}

fn scan_tetragrams() -> HashMap<Vec<u8>, i32> {
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

fn decipher(text: &Vec<u8>, key: [u8; 26]) -> Vec<u8> {
    let mut deciphered: Vec<u8> = Vec::new();
    for i in 0..text.len() {
        for j in 0..26 {
            if LETTERS[j] == text[i] {
                deciphered.push(key[j])
            }
        }
    }
    deciphered
}

fn read(path: &str) -> Vec<u8> {
    let f: Result<Vec<u8>, std::io::Error> = fs::read(path);
    match f {
        Ok(s) => s,
        Err(_) => panic!("Couldn't find file"),
    }
}

fn score(text: &Vec<u8>, tetragrams: &HashMap<Vec<u8>, i32>) -> i32 {
    let mut total: i32 = 0;
    for i in 0..text.len() - 3 {
        let v: Vec<u8> = text[i..i + 4].to_vec();
        if tetragrams.contains_key(&v) {
            total += tetragrams[&v];
        }
    }
    total
}

fn randomise() -> [u8; 26] {
    let mut n: Vec<u8> = LETTERS.to_vec();
    let mut rng = StepRng::new(0, 1);
    let mut irs = Irs::default();
    let _ = irs.shuffle(&mut n, &mut rng);
    let t: Result<[u8; 26], _> = n.try_into();
    match t {
        Ok(arr) => arr,
        Err(_) => panic!("Failed to randomise grid"),
    }
}

fn mutate(k: [u8; 26]) -> [u8; 26] {
    let mut n = k;
    let mut rng = rand::thread_rng();
    let i1 = rng.gen_range(0..26);
    let mut i2 = rng.gen_range(0..26);
    while i1 == i2 {
        i2 = rng.gen_range(0..25);
    }
    let temp = n[i1];
    n[i1] = n[i2];
    n[i2] = temp;
    n
}

fn acceptance_probability(delta_e: f32, temp: f32) -> f32 {
    (1.0 / delta_e) * temp
}

fn simulated_annealing(ciphertext: &Vec<u8>, max_constant: i32, max_temp: f32, k: f32) -> String {
    let mut temp: f32 = max_temp;
    let tetragrams = scan_tetragrams();
    let mut constant: i32 = 0;
    let mut old_key: [u8; 26] = randomise();
    while constant < max_constant {
        let new_key = mutate(old_key);
        let s_old: f32 = 1.0;
        let s_new: f32 = score(&decipher(ciphertext, new_key), &tetragrams) as f32
            / score(&decipher(ciphertext, old_key), &tetragrams) as f32;
        let delta_e = s_old / s_new;
        if delta_e < 1.0 {
            //new key is better
            old_key = new_key;
            constant = 0;
        } else {
            let mut rng = rand::thread_rng();
            let x: f32 = rng.gen();
            let p = acceptance_probability(delta_e, temp);
            if x <= p {
                old_key = new_key;
                constant = 0;
            } else {
                constant += 1;
            }
        }
        temp *= k;
    }
    String::from_utf8(decipher(ciphertext, old_key)).unwrap()
}

fn main() {
    let original = read("src/ciphertext.txt");
    let ciphertext = format_text(original);
    println!("{}", simulated_annealing(&ciphertext, 1000, 1.0, 0.95));
}
