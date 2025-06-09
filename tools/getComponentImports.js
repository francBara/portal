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
        CallExpression(path) {
            const callee = path.get("callee");
            if (callee.isIdentifier({ name: "require" })) {
              const arg = path.get("arguments")[0];
              if (arg.isStringLiteral()) {
                result.imports.push(arg.node.value);
              }
            }
          },
    });

    console.log(JSON.stringify(result));
});