import { AnnotationType, InsertAnnotation } from "@shared/schema";

// Regex patterns for different annotation types
const ANNOTATION_PATTERNS = {
    // Match @cbText variableName: "value"
    text: /@cb(Text)\s+(\w+)(?::\s*["']([^"']*)["'])?/,

    // Match @cbNumber variableName: 123
    number: /@cb(Number)\s+(\w+)(?::\s*(-?\d+(?:\.\d+)?))?/,

    // Match @cbToggle variableName: true/false
    boolean: /@cb(Toggle)\s+(\w+)(?::\s*(true|false))?/,

    // Match @cbColor variableName: "#color"
    color: /@cb(Color)\s+(\w+)(?::\s*["']?(#[a-fA-F0-9]{3,6})["']?)?/,

    // Match @cbSelect variableName: "option", // [option1, option2]
    select: /@cb(Select)\s+(\w+)(?::\s*["']([^"']*)["'])?(?:,?\s*\/\/\s*\[(.*?)\])?/,

    // Match @cbMultiSelect variableName: ["option1", "option2"], // [option1, option2, option3]
    multiSelect: /@cb(MultiSelect)\s+(\w+)(?::\s*\[(.*?)\])?(?:,?\s*\/\/\s*\[(.*?)\])?/,

    // Match @cbComponent variableName: { ... }
    component: /@cb(Component)\s+(\w+)(?::\s*({.*?}|{[^{}]*(?:{[^{}]*}[^{}]*)*}))?/s,
};

// Parse a single line for annotations
export function parseAnnotation(
    line: string,
    lineNumber: number,
    fileId: number
): InsertAnnotation | null {
    for (const [type, pattern] of Object.entries(ANNOTATION_PATTERNS)) {
        const match = line.match(pattern);

        if (match) {
            // Extract marker, name, and value from the match
            const marker = `@cb${match[1]}`;
            const name = match[2];
            let defaultValue: string | undefined;
            let options: any | undefined;

            switch (type) {
                case 'text':
                case 'color':
                    defaultValue = match[3] || '';
                    break;
                case 'number':
                    defaultValue = match[3] || '0';
                    break;
                case 'boolean':
                    defaultValue = match[3] || 'false';
                    break;
                case 'select':
                    defaultValue = match[3] || '';
                    if (match[4]) {
                        try {
                            // Parse the options array from the comment
                            options = match[4].split(',').map(o => o.trim().replace(/["']/g, ''));
                        } catch (e) {
                            console.error(`Error parsing options for select annotation: ${e}`);
                        }
                    }
                    break;
                case 'multiSelect':
                    try {
                        if (match[3]) {
                            // Parse the default selections from the array notation
                            defaultValue = match[3].replace(/["']/g, '').split(',').map(o => o.trim()).join(',');
                        } else {
                            defaultValue = '';
                        }

                        if (match[4]) {
                            // Parse the options array from the comment
                            options = match[4].split(',').map(o => o.trim().replace(/["']/g, ''));
                        }
                    } catch (e) {
                        console.error(`Error parsing multiSelect annotation: ${e}`);
                    }
                    break;
                case 'component':
                    try {
                        if (match[3]) {
                            // For components, store the whole object as a string to be parsed later
                            defaultValue = match[3];
                        }
                    } catch (e) {
                        console.error(`Error parsing component annotation: ${e}`);
                    }
                    break;
            }

            return {
                fileId,
                marker,
                name,
                type: type as AnnotationType,
                line: lineNumber,
                defaultValue: defaultValue ?? null,
                options: options ?? null,
                description: null,
            };
        }
    }

    return null;
}

// Parse a file for annotations
export function parseFile(fileContent: string, fileId: number): InsertAnnotation[] {
    const lines = fileContent.split('\n');
    const annotations: InsertAnnotation[] = [];

    lines.forEach((line, index) => {
        const lineNumber = index + 1;
        const annotation = parseAnnotation(line.trim(), lineNumber, fileId);

        if (annotation) {
            annotations.push(annotation);
        }
    });

    return annotations;
}
