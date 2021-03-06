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


Example
-------

We all know the small UNIX tool “uptime”:

.. code-block:: sh

   $ uptime
    17:25:40 up 20:39,  3 users,  load average: 0.02, 0.07, 0.10

If you need the 1-minute load average, you can write:

.. code-block:: sh

   $ jfy uptime | jq .load1
   0.02

(Yes, ``cut -d " " -f1 /proc/loadavg`` would also work in this case.  But I
needed a really simple tool for starting.)


Status of this project
----------------------

This project is in its early infancy.  I have only wrapped “uptime” because it
is so simple.  Besides that, I just created the common code which works as a
dispatch for the wrappers.

I had the idea and wanted not to forget it.  Therefore, I started this
repository.  When I have some time, I will create more wrappers.  If you are
interested, get in touch with me!

My personal priority is the Btrfs command line tool ``btrfs``.  It has
extensive output but no ``--json`` option.  ``colormgr`` is another candidate.
And there are many others.


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

You may influence jfy’s behaviour by setting the environment variable
``JFY_SETTINGS`` to a JSON string, like this:

.. code-block:: sh

  JFY_SETTINGS='{"exitCode": 1}' jfy ls -l

Alternatively:

.. code-block:: sh

  export JFY_SETTINGS='{"exitCode": 1}'
  jfy ls -l


exitCode
........

The default exit code of jfy (the one used if jfy encounters an error rather
than the tool it wraps) is an arbitrary 221.  By using the setting “exitCode”,
you can choose a code that is distinguishable from the code returned by the
wrapped program.


version
.......

By default, jfy uses the latest version of the respective output transformer.
But the transformation code may be under development and its coverage of the
respective wrapped program’s output may improve.  At the same time, your
scripts may depend on the former JSON data structure.  Therefore, you can
request for a certain version of the JSON structure.

This version can be passed in the “version” setting.  Its format is
``YYYYMMDD``, i.e., a date.  Just pass the date version of the jfy that you are
using, and you should be safe.  Its default value is the special number
99999999, which means to take the latest version.


How to contribute a wrapper
---------------------------

You create a directory with the name of the program to wrap.  It must contain a
Go package of the same name, with at least a ``Handle`` function.  Its
signature is:

.. code-block:: go

   func Handle(settings lib.SettingsType, stdout, stderr []byte,
               args ...string) (output any, errors any, err error) {

Add this handler to the ``dispatchers`` map in the main program.  The
parameters and return values mean the following:

``settings``
  The global settings, like “ExitCode” and “Version” (see `Settings`_).

``stdout``
  The raw stdout content of the wrapped program.

``stderr``
  The raw stderr content of the wrapped program.

``args``
  The command line arguments passed to the wrapped program.  For example, if
  jfy was invoked with ``jfy ls -l``, ``args`` contains only ``{"ls", "-l"}``.

``output``
  An arbitrary Go data structure (well, it must be JSON-serialisable) that
  represents stdout of the wrapped program completely.  Most often, it is of
  type ``map[string]any``.

``errors``
  An arbitrary Go data structure (it must be JSON-serialisable) that represents
  stderr of the wrapped program completely.  Most often, it is of type
  ``map[string]any``.  It should be ``nil`` if there is no error output.

``err``
  Any errors that occurred along the way.  Note that your handler must not
  terminate jfy, nor must it send anything to stdout or stdin.

Your directory should contain a README explaining the structures of ``output``
and ``errors``.

Your wrapper must return an error if it is not sure doing the right thing.  In particular,

- variants of the wrapped program not handled by you (e.g. of another operating
  system),
- command line options affecting the output in a way not handled by you,
- command line options unknown to you,
- environment variables affecting the output in a way not handled by you
  (e.g. ``LANG``)

must make your code return an error.  Such checks can never be perfect but try
to be as careful as possible.

When you are finished, you are invited to create a PR for your addition.  Mind
the AGPL this program is published under.  Thank you!


..  LocalWords:  jfy Stdout stderr JSONification cp mv jq ls’ Stdin jfy’s
..  LocalWords:  stdout exitCode
