use std::collections::HashMap;
use std::fs;
use std::fs::File;
use std::io::{BufRead, BufReader};

fn read(path: &str) -> Vec<u8> {
    let f: Result<Vec<u8>, std::io::Error> = fs::read(path);
    match f {
        Ok(s) => s,
        Err(_) => panic!("Couldn't find file"),
    }
}

fn format_text(original: Vec<u8>) -> Vec<u8> {
    let mut result: Vec<u8> = Vec::new();
    for b in original {
        if (65..=90).contains(&b) {
            result.push(b)
        }
    }
    result
}

fn decipher(text: Vec<u8>, key: u8) -> Vec<u8> {
    let mut deciphered: Vec<u8> = Vec::new();
    for b in text {
        let d = (b - 65 + key) % 26 + 65;
        deciphered.push(d)
    }
    deciphered
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

fn gen_reference(n: usize) -> Vec<u8> {
    let r: Vec<u8> = format_text(read("../referenceText.txt"));
    let mut ref_text: Vec<u8> = Vec::new();
    for i in 0..n {
        ref_text.push(r[i]);
    }
    ref_text
}

fn bruteforce(text: Vec<u8>) {
    let tetragrams = scan_tetragrams();
    let ref_score = score(&gen_reference(text.len()), &tetragrams);
    for i in 0..25 {
        let p = decipher(text.clone(), i);
        if score(&p, &tetragrams) * 10 >= ref_score * 8 {
            println!("{}", String::from_utf8(p).unwrap());
            return;
        }
    }
    panic!("Found no key that works, unlikely to be Caesar cipher")
}

fn main() {
    let original = read("src/ciphertext.txt");
    let ciphertext = format_text(original);
    bruteforce(ciphertext);
}
