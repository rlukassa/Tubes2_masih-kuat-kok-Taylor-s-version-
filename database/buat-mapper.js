import axios from 'axios';
import * as cheerio from 'cheerio';
import fs from 'fs';

const url = 'https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)';
const response = await axios.get(url);
const $ = cheerio.load(response.data);

let seen = new Set(); // Untuk elemen unik
let mapper = {};      // Untuk hasil akhir JSON

// mengkapitalkan huruf pertama dari string
const capitalize = (str) => str.charAt(0).toUpperCase() + str.slice(1);

$('img').each((i, img) => { // mencari tag img
  const dataSrc = $(img).attr('data-src'); // mendapatkan atribut data-src
  if (dataSrc && dataSrc.includes('little-alchemy/images') && dataSrc.includes('svg')) { // filter kalo ada di little-alchemy/images dan ada svg
    const match = dataSrc.match(/\/([^/]+)_\d+\.svg/); // di ambil nama file dari data-src
    if (match) {
      const elementName = capitalize(match[1].toLowerCase()); // mengkapitalkan huruf pertama dari nama file
      if (!seen.has(elementName)) {
        seen.add(elementName); // menambah elemen ke set
        const cleanUrl = dataSrc.split('/revision/')[0]; // Memotong URL setelah "/revision/"
        mapper[elementName] = cleanUrl; // Menambah elemen ke mapper
      }
    }
  }
});

fs.writeFileSync('mapper2.json', JSON.stringify(mapper, null, 2), 'utf8');
console.log(`Saved ${Object.keys(mapper).length} elements to mapper2.json`);