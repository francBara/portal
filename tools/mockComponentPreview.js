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

    const result = {};

    traverse(ast, {
        FunctionDeclaration(path) {
            if (path.node.leadingComments) {
                for (let comment of path.node.leadingComments) {
                    if (comment.value.trim().startsWith("@portal")) {
                        result.componentName = path.node.id.name;
                    }
                }
            }
        }
    });

    //result.sourceCode = recast.print(ast).code;

    console.log(JSON.stringify(result));
});