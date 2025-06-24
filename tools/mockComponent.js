const traverse = require('@babel/traverse').default;
const recast = require('recast');

let input = "";
process.stdin.on('data', chunk => input += chunk);

process.stdin.on("end", () => {
    input = JSON.parse(input);

    const ast = recast.parse(input.sourceCode, {
        parser: require("recast/parsers/babel-ts"),
    });

    console.error(input.mocks);

    traverse(ast, {
        VariableDeclaration(path) {
            if (!path.node.leadingComments) {
                return;
            }

            isMock = false;

            for (let comment of path.node.leadingComments) {
                if (comment.value.includes("@portal") && comment.value.includes("mock")) {
                    isMock = true;
                    break
                }
            }

            if (!isMock) {
                return;
            }

            for (let decl of path.node.declarations) {
                if (decl.id.name in input.mocks) {
                    decl.init = recast.types.builders.identifier(input.mocks[decl.id.name]);
                }
            }
        }
    });
    
    console.log(recast.print(ast).code);
});