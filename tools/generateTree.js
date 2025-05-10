const fs = require('fs');
const parser = require('@babel/parser');
const traverse = require('@babel/traverse').default;
const t = require('@babel/types');

const sourceCode = fs.readFileSync(process.argv[2], 'utf8');

const ast = parser.parse(sourceCode, {
    sourceType: 'module',
    plugins: ['typescript', 'jsx'],
});

const components = {};

traverse(ast, {
    FunctionDeclaration(path) {
        const rootName = path.node.id.name;

        for (let el of path.node.body.body) {
            if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                components[rootName] = collectJSX(el.argument);
            }
        }
    },
    VariableDeclarator(path) {
        const rootName = path.node.id.name;

        if (path.node.init && path.node.init.type === "ArrowFunctionExpression") {
            for (let el of path.node.init.body.body) {
                if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                    components[rootName] = collectJSX(el.argument);
                }
            }
        }
    }
});

function collectJSX(node) {
    if (!t.isJSXElement(node)) return;

    const element = {
        type: node.openingElement.name.name || 'Unknown',
        row: node.openingElement.loc.start.line,
        properties: [],
        children: [],
    };

    node.openingElement.attributes.forEach(attr => {
        if (t.isJSXAttribute(attr) && t.isJSXIdentifier(attr.name) && attr.name.name === "className") {
            if (t.isStringLiteral(attr.value)) {
                element.row = attr.loc.start.line;
                element.properties = parseTailwindString(attr.value.value);
            }
            else if (t.isJSXExpressionContainer(attr.value) && t.isTemplateLiteral(attr.value.expression)) {
                if (attr.value.expression.expressions.length > 0) {
                    //TODO: Implement expression container parsing
                }
            }
        }
    });

    node.children.forEach(child => {
        if (t.isJSXElement(child)) {
            element.children.push(collectJSX(child));
        } else if (t.isJSXText(child)) {
            const trimmed = child.value.trim();
            if (trimmed) {
                element.children.push({ type: 'text', value: trimmed });
            }
        }
    });

    return element;
}

function parseTailwindString(tailwind) {
    const splitted = tailwind.split(" ");

    const result = [];

    for (let word of splitted) {
        const lastHyphen = word.lastIndexOf("-");

        let prefix = word.substring(0, lastHyphen);
        let value = word.substring(lastHyphen + 1);

        if (prefix.length === 0) {
            prefix = value;
            value = "";
        }

        result.push({ prefix, value });
    }

    return result;
}

console.log(JSON.stringify(components, null, 2));
