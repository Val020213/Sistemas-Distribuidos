module.exports = {
  // Type check TypeScript files
  '**/*.(ts|tsx)': () => 'npm run tsc --noEmit',

  // Lint then format TypeScript and JavaScript files
  '**/*.(ts|tsx|js)': (filenames) => [
    `npm run eslint --fix ${filenames.join(' ')}`,
    $(`npm run format --write ${filenames.join(' ')}`),
  ],

  // Format MarkDown and JSON
  '**/*.(md|json)': (filenames) => $(`npm run format --write ${filenames.join(' ')}`),
};