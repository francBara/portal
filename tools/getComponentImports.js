const traverse = require('@babel/traverse').default;
const recast = require('recast');

let input = "";
process.stdin.on('data', chunk => input += chunk);

process.stdin.on("end", () => {
    input = JSON.parse(input);

    const sourceCode = input.sourceCode;

    const ast = recast.parse(sourceCode, {
        parser: require("recast/parsers/babel-ts")
    });

    const result = {imports: []};

    traverse(ast, {
        ImportDeclaration(path) {
            result.imports.push(path.node.source.value);
        },
    });

    console.log(JSON.stringify(result));
});