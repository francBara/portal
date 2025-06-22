const traverse = require('@babel/traverse').default;
const { default: generate } = require('@babel/generator');
const recast = require('recast');
const t = require('@babel/types');

let input = "";
process.stdin.on('data', chunk => input += chunk);

process.stdin.on("end", () => {
    input = JSON.parse(input);

    const ast = recast.parse(input.sourceCode, {
        parser: require("recast/parsers/babel-ts")
    });

    traverse(ast, {
        ObjectProperty(path) {
            if (t.isIdentifier(path.node.key, { name: 'content' }) && t.isArrayExpression(path.node.value)) {
                const existing = path.node.value.elements.map(e => e.value);
                if (!existing.includes(input.newPath)) {
                    path.node.value.elements.push(t.stringLiteral(input.newPath));
                }
            }
        }
    });

    console.log(recast.print(ast).code);
});