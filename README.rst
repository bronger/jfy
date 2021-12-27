jfy – JSONification of program output
=====================================

The UNIX idea of passing plain data from one command-line program to the next
through pipes is ingenious.  It makes simple things simple and complex things
possible.  However, sometimes you may want also the complex things to be
simple.  Here, jfy comes into play.

jfy is a wrapper around common UNIX tools like ls, cp, mv etc. and converts
their output to JSON.  This JSON can then be piped into `jq`_, making it easy
to extract what you want safely and conveniently.

.. _jq: https://stedolan.github.io/jq/


Usage
-----

Call

.. code-block:: sh

  jfy ls -l

to get ls’ long output as JSON.

Signals sent to jfy are passed to the child process.  Stdin is passed to the
child as well.  Stdout and stderr of the child are converted to JSON and sent
to jfy’s stdout and stderr, respectively.  The exit code of the child is
returned as jfy’s exit code.


Settings
--------

You may influence jfy’s behaviour by settings the environment variable
``JFY_SETTINGS`` to a JSON string, like this:

.. code-block:: sh

  JFY_SETTINGS='{"exitCode": 1}' jfy ls -l

Alternatively:

.. code-block:: sh

  export JFY_SETTINGS='{"exitCode": 1}'
  jfy ls -l

Currently, only one settings exists, namely “exitCode”.  The default exit code
of jfy (the one used if jfy encounters an error rather than the tool it wraps)
is an arbitrary 221.  By using the setting “exitCode”, you can choose a code
that is distinguishable from the code returned by the wrapped program.


..  LocalWords:  jfy
