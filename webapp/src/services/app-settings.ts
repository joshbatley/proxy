const apibase = ' http://localhost:5000/admin';

export default ({
  selector: `${apibase}/collections/selector`,
  response: `${apibase}/responses`,
  endpoint: `${apibase}/endpoints`,
});
