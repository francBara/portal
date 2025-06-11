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
                    
                    //TODO: Support whitespaces, make better regex
                    //TODO: Check that all props are mocked
                    if (comment.startsWith("@portal")) {
                        result.componentName = path.node.id.name;

                        const mock = comment.replace(/^(@portal)/, "");

                        result.mock = JSON.parse(mock);
                    }

                    break;
                }
            }
        }
    });

    //result.sourceCode = recast.print(ast).code;

    console.log(JSON.stringify(result));
});