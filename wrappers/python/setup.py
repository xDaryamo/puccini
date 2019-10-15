#!/usr/bin/env python

import os.path, setuptools, subprocess, tempfile, shutil

with open(os.path.join(os.path.dirname(__file__), 'README.rst')) as f:
    readme = f.read()

script = '''\
ROOT="{}"

# Install Go
curl https://storage.googleapis.com/golang/go1.13.1.linux-amd64.tar.gz --silent --location | tar -xz
export PATH="$PATH:go/bin"

# Fetch repository
REPO=puccini
git clone --depth 1 https://github.com/tliron/puccini "$REPO"

# Build library
"$REPO/scripts/build-library.sh"
mv "$REPO/dist/libpuccini.so" "$ROOT/puccini/" 
'''.format(os.path.abspath(os.path.dirname(__file__)).replace('"', '\\"'))

t = tempfile.mkdtemp()
try:
    subprocess.check_output(('bash',  '-o', 'pipefail', '-euxc', script), cwd=t)
finally:
    shutil.rmtree(t)

setuptools.setup(
    name='puccini',
    version='0.1',
    description='Puccini',
    long_description=readme,
    license='Apache License 2.0',
    author='Tal Liron',
    author_email='tal.liron@gmail.com',
    url='https://github.com/tliron/puccini',
    download_url='https://github.com/tliron/puccini/releases',
    classifiers=[
        'Development Status :: 4 - Beta',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache Software License'],

    packages=['puccini'],
    package_dir={'puccini': 'puccini'},
    package_data={'puccini': ['libpuccini.so']},
    install_requires=['ruamel.yaml'])
