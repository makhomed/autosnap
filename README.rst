========================
autosnap (version 2.3.0)
========================

ZFS snapshot automation tool

Installation
------------

- ``cd /opt``
- ``git clone https://github.com/makhomed/autosnap.git autosnap``

Also you need to install python3:

.. code-block:: none

    # yum install python3

Upgrade
-------

- ``cd /opt/autosnap``
- ``git pull``

Configuration
-------------

- ``vim /opt/autosnap/autosnap.conf``
- write to config something like this:

.. code-block:: none

    interval hourly 24
    interval daily  15
    interval weekly 8

    exclude **-dump
    exclude **-swap

    ignore tank/backup/**
    ignore tank/mirror/**


Configuration file allow comments, from symbol ``#`` to end of line.

Configuration file has only four directives:
``interval``, ``exclude``, ``include`` and ``ignore``.

Syntax of interval directive: ``interval <name> <count>``.
``<name>`` is name of interval, must be unique.
``<count>`` is count of snapshots to save for interval ``<name>``.

Interval name can be ``frequent``, ``hourly``, ``daily``, ``weekly``, ``monthly``, ``yearly`` or something else.
Each interval name should be configured via cron for running.

Syntax of ``include``, ``exclude`` and ``ignore`` directives are the same:
``exclude <pattern>`` or ``include <pattern>`` or ``ignore <pattern>``.

By default all datasets are included. But you can exclude some datasets
by name or by pattern. Pattern is rsync-like, ``?`` means any one symbol,
``*`` means any symbols except ``/`` symbol, ``**`` means any symbols.

First match win, and if it was directive ``exclude`` - dataset will be excluded,
if it was directive ``include`` - dataset will be included.

``exclude``, ``include`` and ``ignore`` directives allowed only on global level.

``ignore`` command works as ``exclude`` command, and also ``ignore`` datasets will be ignored in the ``autosnap -x`` command output.

``interval`` directive allowed to be configured for each dataset separately.
For example:

.. code-block:: none

    interval hourly 24
    interval daily  15
    interval weekly 8

    [tank/kvm-stage-elastic]

    interval hourly 24
    interval daily  7
    interval weekly 4

    [tank/kvm-stage-mysqld]

    interval hourly 24
    interval daily  7
    interval weekly 4

Schedule autosnap
-----------------

- ``vim /etc/cron.d/autosnap``
- write to cron file something like this:

.. code-block:: none

    0 * * * * root /opt/autosnap/autosnap -x hourly
    0 0 * * * root /opt/autosnap/autosnap -x daily
    0 0 * * 7 root /opt/autosnap/autosnap -x weekly

Command line
------------

At start ``autosnap`` read config from ``/opt/autosnap/autosnap.conf``.

One and only one command must be specified in command line. This command must
be the name of interval from configuration file.

During execution, autosnap will create one new snapshot for each included dataset
and delete all oldest snapshots exceeding the allowed snapshots count for given interval.

If autosnap started without any parameters - then help printed and program will exit.

If autosnap started with ``-x`` parameter - it will list all snapshots, not managed by autosnap.
The ``-x`` parameter can be combined with interval name.

If autosnap started with ``-d`` parameter - it will delete all expired snapshots, managed by autosnap.
The ``-d`` parameter can not be combined with interval name.

