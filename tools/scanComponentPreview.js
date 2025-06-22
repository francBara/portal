const traverse = require('@babel/traverse').default;
const recast = require('recast');

function parseKeyValues(input) {
  const pairs = input.split(/\s+/);
  const result = {};
  
  pairs.forEach(pair => {
    const [key, value] = pair.split('=');
    if (key && value !== undefined) {
      result[key] = isNaN(value) ? value : Number(value);
    }
  });

  return result;
}

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
                    if (comment.includes("@portal")) {
                        // Get component name
                        result.componentName = path.node.id.name;

                        const arguments = parseKeyValues(comment);

                        if (arguments.h) {
                            result.boxHeight = arguments.h;
                        }
                        if (arguments.w) {
                            result.boxWidth = arguments.w;
                        }

                        // Get props mock
                        const propsIndex = comment.indexOf("props");

                        if (propsIndex !== -1) {
                            result.mock = JSON.parse(comment.slice(propsIndex +  5));
                        }
                        
                        break;
                    }
                }
            }
        }
    });

    //result.sourceCode = recast.print(ast).code;

    console.log(JSON.stringify(result));
});