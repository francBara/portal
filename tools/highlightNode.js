const traverse = require('@babel/traverse').default;
const t = require('@babel/types');
const recast = require('recast');

const highlightClass = "outline outline-2 outline-red-500 outline-offset-2";

const currentId = 0;

function updateNode(node, highlightedNodeId) {
    if (!t.isJSXElement(node)) return;

    for (let attr of node.openingElement.attributes) {
        if (t.isJSXAttribute(attr) && t.isJSXIdentifier(attr.name) && attr.name.name === "className") {
            if (t.isStringLiteral(attr.value)) {
                if (currentId === highlightedNodeId) {
                    if (!attr.value.value.includes(highlightClass)) {
                        attr.value.value += " " + highlightClass;
                    }
                }
                else {
                    attr.value.value = attr.value.value.replace(highlightClass, "");
                }
            }
            else if (t.isJSXExpressionContainer(attr.value) && t.isTemplateLiteral(attr.value.expression)) {
                //TODO: Implement expression container parsing
            }
        }
    }
    
    currentId++;

    for (let i = 0; i < node.children.length; i++) {
        if (t.isJSXElement(node.children[i])) {
            updateNode(node.children[i], highlightedNodeId);
        }
    }
}

let input = "";
process.stdin.on('data', chunk => input += chunk);

process.stdin.on("end", () => {
    input = JSON.parse(input);

    const ast = recast.parse(input.sourceCode, {
        parser: require("recast/parsers/babel-ts")
    });

    traverse(ast, {
        FunctionDeclaration(path) {
            const rootName = path.node.id.name;

            if (rootName !== input.rootName) {
                return;
            }

            for (let el of path.node.body.body) {
                if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                    updateNode(el.argument, input.nodeId);
                }
            }
        },
        VariableDeclarator(path) {
            const rootName = path.node.id.name;

            if (rootName !== input.rootName) {
                return;
            }

            if (path.node.init && path.node.init.type === "ArrowFunctionExpression") {
                for (let el of path.node.init.body.body) {
                    if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                        updateNode(el.argument, input.nodeId);
                    }
                }
            }
        }
    });

    console.log(recast.print(ast).code);
});