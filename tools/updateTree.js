const fs = require('fs');
const traverse = require('@babel/traverse').default;
const t = require('@babel/types');
const recast = require('recast');

function updateNode(node, newValue) {
    if (!t.isJSXElement(node)) return;

    for (let attr of node.openingElement.attributes) {
        if (t.isJSXAttribute(attr) && t.isJSXIdentifier(attr.name) && attr.name.name === "className") {
            if (t.isStringLiteral(attr.value)) {
                console.error(`Updating ${attr.value.value} with ${parseTailwind(newValue.properties)}`);
                attr.value.value = parseTailwind(newValue.properties);
            }
            else if (t.isJSXExpressionContainer(attr.value) && t.isTemplateLiteral(attr.value.expression)) {
                //TODO: Implement expression container parsing
            }
        }
    }

    let newValueIndex = 0;

    for (let i = 0; i < node.children.length; i++) {
        if (t.isJSXElement(node.children[i])) {
            if (newValue.children[newValueIndex]) {
                updateNode(node.children[i], newValue.children[newValueIndex]);
            }
            newValueIndex += 1;
        }
    }
}

function parseTailwind(properties) {
    let tailwindString = "";

    for (let p of properties) {
        if (p.value.length == 0) {
            tailwindString += p.prefix;
        }
        else {
            tailwindString += p.prefix + "-" + p.value;
        }

        tailwindString += " ";
    }

    return tailwindString.trim();
}

let input = "";
process.stdin.on('data', chunk => input += chunk);

process.stdin.on("end", () => {
    input = JSON.parse(input);

    const sourceCode = input.sourceCode;

    const ast = recast.parse(sourceCode, {
        parser: require("recast/parsers/babel-ts")
    });

    const components = input.components;

    traverse(ast, {
        FunctionDeclaration(path) {
            const rootName = path.node.id.name;

            for (let el of path.node.body.body) {
                if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                    updateNode(el.argument, components[rootName]);
                }
            }
        },
        VariableDeclarator(path) {
            const rootName = path.node.id.name;

            if (path.node.init && path.node.init.type === "ArrowFunctionExpression") {
                for (let el of path.node.init.body.body) {
                    if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                        updateNode(el.argument, components[rootName]);
                    }
                }
            }
        }
    });

    console.log(recast.print(ast).code);
});