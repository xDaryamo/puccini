// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/compare_version.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.2.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.2.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.2

function compare(a, b) {
	if ((a.$type !== 'version') || (b.$type !== 'version'))
		throw 'both values must be of type "version"';
	if (a.major !== b.major)
		return a.major < b.major ? -1 : 1;
	if (a.minor !== b.minor)
		return a.minor < b.minor ? -1 : 1;
	if (a.fix !== b.fix)
		return a.fix < b.fix ? -1 : 1;
	var aq = a.qualifier.toLowerCase();
	var bq = b.qualifier.toLowerCase();
	if (aq !== bq) // note: the qualifier is compared alphabetically, *not* semantically
		return aq < bq ? -1 : 1;
	if (a.build !== b.build)
		return a.build < b.build ? -1 : 1;
	return 0;
}
`
}
