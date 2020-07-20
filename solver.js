const solveCaptcha = require('./solver/index.js');

(async () => {
    try {
      const response = await solveCaptcha(process.argv[2], base64decode(process.argv[3]));
      console.log(response);
    } catch (error) {
      console.log(error);
    }
})();

function base64decode(base64text){
    return Buffer.from(base64text, 'base64').toString('utf8');
}