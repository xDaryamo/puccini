#include "cloud_puccini_TOSCA.h"
#include "libpuccini.h"
#include <stdlib.h>

JNIEXPORT jstring JNICALL Java_cloud_puccini_TOSCA__1Compile
  (JNIEnv *env, jclass cls, jstring url, jstring inputs, jstring quirks)
{
	const char *url_ = (*env)->GetStringUTFChars(env, url, 0);
	const char *inputs_ = (*env)->GetStringUTFChars(env, inputs, 0);
	const char *quirks_ = (*env)->GetStringUTFChars(env, quirks, 0);

	char *result = Compile((char *) url_, (char *) inputs_, (char *) quirks_);

	(*env)->ReleaseStringUTFChars(env, url, url_);
	(*env)->ReleaseStringUTFChars(env, inputs, inputs_);
	(*env)->ReleaseStringUTFChars(env, quirks, quirks_);

	jstring result_ = (*env)->NewStringUTF(env, result);
	free(result);
	return result_;
}
