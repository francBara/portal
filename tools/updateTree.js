const fs = require('fs');
const traverse = require('@babel/traverse').default;
const t = require('@babel/types');
const recast = require('recast');

currentId = 0;

function updateNode(node, newValue) {
    if (t.isJSXText(node)) {
        if (newValue.type !== "text") {
            throw("Input tree not consistent, expected text, got " + newValue.type);
        }
        node.value = newValue.properties[0].value;
        currentId++;
        return;
    }

    if (!t.isJSXElement(node)) return;

    if (node.openingElement.name.name !== newValue.type) {
        throw("Input tree not consistent, expected " + node.openingElement.name.name + ", got " + newValue.type);
    }

    for (let attr of node.openingElement.attributes) {
        if (t.isJSXAttribute(attr) && t.isJSXIdentifier(attr.name) && attr.name.name === "className") {
            if (t.isStringLiteral(attr.value)) {
                newTwString = parseTailwind(newValue.properties);

                attr.value.value = newTwString;
            }
            else if (t.isJSXExpressionContainer(attr.value) && t.isTemplateLiteral(attr.value.expression)) {
                //TODO: Implement expression container parsing
            }
        }
    }

    currentId++;

    let newValueIndex = 0;

    for (let i = 0; i < node.children.length; i++) {
        if (t.isJSXElement(node.children[i]) || (t.isJSXText(node.children[i]) && node.children[i].value.trim())) {
            updateNode(node.children[i], newValue.children[newValueIndex]);
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

            if (!components[rootName] || !path.node.leadingComments) {
                return;
            }

            for (let el of path.node.body.body) {
                if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                    updateNode(el.argument, components[rootName]);
                }
            }
        },
        VariableDeclarator(path) {
            const rootName = path.node.id.name;

            if (!components[rootName] || !path.node.leadingComments) {
                return;
            }

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