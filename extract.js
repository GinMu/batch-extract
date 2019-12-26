const { execSync } = require("child_process");
const program = require("commander");
const fs = require("fs");
const path = require("path");

program
  .usage("[options] <file ...>")
  .option("-s, --source <path>", "源目录")
  .option("-d, --destination <path>", "目标目录");

program.parse(process.argv);

const { source, destination } = program;

fs.readdir(source, "utf-8", function callback(err, files) {
  if (err) {
    return console.log(err);
  }

  const len = files.length;

  for (let i = 0; i < len; i++) {
    const file = files[i];
    const absPath = source + file;
    const extname = path.extname(absPath);
    if (extname !== ".rar" && extname !== ".zip") {
      continue;
    }
    // 因为文件名包含三位数目录和其他汉字，所以获取三位数目录, 作为解压后目录
    const rename = file.match(/\d{3}/g)[0];
    const filename = rename + extname;
    const destFile = `${destination + filename}`;
    try {
      let extractCmd;
      const destFolder = `${destination + rename}`;
      if (extname === ".rar") {
        extractCmd = `unrar e ${destFile} ${destFolder}`;
      } else {
        extractCmd = `ditto -x -k ${destFile} ${destFolder}`;
      }

      fs.copyFileSync(absPath, destFile);
      execSync(`mkdir ${destFolder}`);
      execSync(extractCmd);
      // const files = fs.readdirSync(`${destFolder}`, "utf-8");
      // if (files.length === 1) {
      //   execSync(`mv -f ${destFolder + "/" + files[0].replace("\\", "\\\\") + "/*"} ${destFolder}`);
      //   execSync(`rm -rf ${destFolder + "/" + files[0].replace("\\", "\\\\")}`);
      // }
    } catch (error) {
      console.log(filename + ": " + error.message);
    } finally {
      execSync(`rm ${destFile}`);
    }
  }
});

// unzip src -d dest
// unrar e src dest
