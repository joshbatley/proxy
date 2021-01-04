const apibase = ' http://localhost:5000/admin';

export default ({
  collections: `${apibase}/collections`,
  response: `${apibase}/responses`,
  endpoint: `${apibase}/endpoints`,
  rules: `${apibase}/rules`,
});
