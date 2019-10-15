#include "puccini_TOSCA.h"
#include "libpuccini.h"
#include <stdlib.h>

JNIEXPORT jstring JNICALL Java_puccini_TOSCA__1Compile
  (JNIEnv *env, jclass cls, jstring url)
{
	const char *url_ = (*env)->GetStringUTFChars(env, url, 0);
	char *clout = Compile((char *) url_);
	(*env)->ReleaseStringUTFChars(env, url, url_);

	jstring r = (*env)->NewStringUTF(env, clout);
	free(clout);
	return r;
}
