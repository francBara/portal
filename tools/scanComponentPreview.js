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
                for (let leadingComment of path.node.leadingComments) {
                    comment = leadingComment.value.trim();
                    
                    if (comment.startsWith("@portal")) {
                        result.componentName = path.node.id.name;
                        result.mock = comment.replace(/^(@portal)/, "");
                    }

                    break;
                }
            }
        }
    });

    //result.sourceCode = recast.print(ast).code;

    console.log(JSON.stringify(result));
});