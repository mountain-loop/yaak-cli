const fs = require("fs");
const path = require("path");

const readme = path.join(__dirname, '..', '..', 'README.md');

fs.copyFileSync(readme, path.join(__dirname, "README.md"))
