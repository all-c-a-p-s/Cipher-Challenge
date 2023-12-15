use rand::rngs::mock::StepRng;
use shuffle::irs::Irs;
use shuffle::shuffler::Shuffler;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fs;
use std::fs::File;
use std::io::{BufRead, BufReader};

#[allow(dead_code)]

struct Key {
    a: [char; 25],
    b: [char; 25],
    c: [char; 25],
    d: [char; 25],
}

struct Coords {
    row: i32,
    column: i32,
}

fn letters() -> [char; 25] {
    let l: [char; 25] = [
        'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
        'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
    ];
    l
}

fn format_text(original: String) -> String {
    let mut result: String = "".to_string();
    for i in 0..original.len() {
        if original.as_bytes()[i] >= 65 && original.as_bytes()[i] <= 90 {
            result.push(original.as_bytes()[i] as char)
        }
    }
    result
}

fn scan_tetragrams() -> HashMap<String, i32> {
    let mut tg: HashMap<String, i32> = HashMap::new();
    let file = File::open("../tetragrams.txt").expect("Failed to load tetragram file");
    let reader = BufReader::new(&file);
    for line in reader.lines() {
        let ln: &str = &line.unwrap_or(String::from("Invalid line"));
        let fields: Vec<&str> = ln.split(", ").collect();
        let n: Result<i32, _> = fields[1].parse();
        let num: i32 = match n {
            Ok(n) => n,
            Err(_) => panic!("Failed to convert tetragram frequence to i32"),
        };
        tg.insert(fields[0].to_string(), num);
    }
    tg
}

fn gen_reference(n: usize) -> String {
    let r: Vec<u8> = format_text(read("../referenceText.txt"))
        .as_bytes()
        .to_vec();
    let mut ref_text: String = String::from("");
    for i in 0..n {
        ref_text = format!("{}{}", ref_text, String::from(r[i] as char));
    }
    ref_text
}

fn scan_dict() -> Vec<String> {
    let words: Vec<String> = read("../challenge_word_list.txt")
        .split(" ")
        .collect::<Vec<&str>>()
        .iter()
        .map(|&s| s.into())
        .collect(); //this was painful lmao
    words[..words.len() - 1].to_vec() //last element is a newline character otherwise
}

fn randomise() -> [char; 25] {
    let mut g: Vec<char> = letters().to_vec();
    let mut rng = StepRng::new(0, 1);
    let mut irs = Irs::default();
    let _ = irs.shuffle(&mut g, &mut rng);
    let t: Result<[char; 25], _> = g.try_into();
    match t {
        Ok(arr) => arr,
        Err(_) => panic!("Failed to randomise grid"),
    }
}

fn read(path: &str) -> String {
    let f: Result<String, std::io::Error> = fs::read_to_string(path);
    match f {
        Ok(s) => s,
        Err(_) => panic!("Couldn't find file"),
    }
}

fn score(text: &String, tetragrams: &HashMap<String, i32>) -> i32 {
    let bytes = text.as_bytes();
    let mut total: i32 = 0;
    for i in 0..text.len() - 3 {
        let v: Vec<u8> = (bytes[i..i + 4]).to_vec();
        let s: Result<String, _> = String::from_utf8(v);
        let slice = match s {
            Ok(s) => s,
            Err(_) => panic!("Failed to parse window in text"),
        };
        if tetragrams.contains_key(&slice) {
            total += tetragrams[&slice];
        }
    }
    total
}

fn fill_grid(keyword: &mut String) -> [char; 25] {
    let mut seen: HashSet<char> = HashSet::new();
    keyword.retain(|c| {
        let s = seen.contains(&c);
        seen.insert(c);
        !s
    }); //remove duplicate letters from keyword

    let last_idx = letters()
        .iter()
        .position(|&x| x == keyword.chars().last().unwrap())
        .unwrap(); //continues filling up grid with next letter after last letter in keyword
    let mut grid: Vec<char> = keyword
        .as_bytes()
        .to_vec()
        .iter()
        .map(|x| *x as char)
        .collect(); //also painful
    for i in last_idx..last_idx + 25 {
        let idx = i % 25;
        if !seen.contains(&letters()[idx]) {
            seen.insert(letters()[idx]);
            grid.push(letters()[idx]);
        }
    }
    let t: Result<[char; 25], _> = grid.try_into();
    match t {
        Ok(arr) => arr,
        Err(_) => panic!("Failed to fill grid"),
    }
}

fn coordinates(grid_index: i32) -> Coords {
    let row: i32 = grid_index / 5;
    let column: i32 = grid_index - row * 5;
    Coords { row, column }
}

fn grid_index(row: i32, column: i32) -> i32 {
    row * 5 + column
}

fn decipher(text: &String, k: Key) -> Option<String> {
    let mut deciphered: String = String::from("");
    let bytes = text.as_bytes();
    for i in (0..text.len() - 1).step_by(2) {
        let v: Vec<u8> = (bytes[i..i + 2]).to_vec();
        let i1: usize = k.b.iter().position(|&x| x == v[0] as char)?;
        let i2: usize = k.c.iter().position(|&x| x == v[1] as char)?;
        let coords1: Coords = coordinates(i1.try_into().ok()?);
        let coords2: Coords = coordinates(i2.try_into().ok()?);

        let d_i1: usize = grid_index(coords1.row, coords2.column)
            .try_into()
            .expect("Failed to get grid index 1"); //coords of deciphered letters
        let d_i2: usize = grid_index(coords2.row, coords1.column)
            .try_into()
            .expect("Failed to get grid index 2");

        let letter1: String = String::from(k.a[d_i1]);
        let letter2: String = String::from(k.d[d_i2]);
        deciphered = format!("{}{}{}", deciphered, letter1, letter2);
    }
    Some(deciphered)
}

fn dictionary_attack(ciphertext: &String) {
    let dict = scan_dict();
    let tetragrams = scan_tetragrams();
    let ref_score = score(&gen_reference(ciphertext.clone().len()), &tetragrams);
    let mut count: i32 = 0;
    for word in &dict {
        println!("{}", count);
        count += 1;
        if word.contains("J") || word.len() == 1 {
            continue;
        }
        for word2 in &dict {
            if word2.contains("J") || word2.len() == 1 {
                continue;
            }
            let k: Key = Key {
                a: letters(),
                b: fill_grid(&mut word.to_string()),
                c: fill_grid(&mut word2.to_string()),
                d: letters(),
            };
            let p = decipher(&ciphertext, k).unwrap();
            if score(&p, &tetragrams) * 10 >= ref_score * 8 {
                println!("{}", p)
            }
        }
    }
}

fn main() {
    let original = read("src/ciphertext.txt");
    let ciphertext: String = format_text(original);
    dictionary_attack(&ciphertext)
}
