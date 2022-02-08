// This scriptlet extracts all artifacts to the output directory

const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

traversal.coerce();

for (let vertexId in clout.vertexes) {
	let vertex = clout.vertexes[vertexId];
	if (!tosca.isNodeTemplate(vertex))
		continue;
	let nodeTemplate = vertex.properties;

	for (let key in nodeTemplate.artifacts) {
		let artifact = nodeTemplate.artifacts[key];

		// If 'puccini.output' is empty, this will be relative to current directory
		let targetPath = puccini.joinFilePath(puccini.output, artifact.filename);

		puccini.log.noticef('extracting "%s" to "%s"', artifact.sourcePath, targetPath);
		puccini.download(artifact.sourcePath, targetPath);

		//puccini.log.noticef('%s', puccini.exec('cat', targetPath));
	}
}
