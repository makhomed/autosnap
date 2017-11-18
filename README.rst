========
autosnap
========

ZFS snapshot automation tool

Installation
------------

 - ``cd /opt``
 - ``git clone https://github.com/makhomed/autosnap.git autosnap```

Configuration
-------------

  - ``vim /opt/autosnap/autosnap.conf``
  - write to config something like this:

.. code-block:: bash

    interval hourly 24
    interval daily  30

    exclude tank
    exclude tank/backup**
    exclude tank/vm

Configuration file allow comments, from symbol ``#`` to end of line.

Configuration file has only three directives:
``interval``, ``exclude`` and ``include``.

Syntax of interval directive: ``interval <name> <count>``.
``<name>`` is name of interval, must be unique.
``<count>`` is count of snapshots to save for interval ``<name>``.

Syntax of ``include`` and ``exclude`` directives are the same:
``exclude <pattern>`` or ``include <pattern>``.

By default all datasets are included. But you can exclude some datasets
by name or by pattern. Pattern is rsync-like, ``?`` means any one symbol,
``*`` means any symbols except ``/`` symbol, ``**`` means any symbols.

First match win, and if it was directive ``exclude`` - dataset will be excluded,
if it was directive ``include`` - dataset will be included.

Schedule autosnap
-----------------

  - ``vim /etc/cron.d/autosnap``
  - write to cron file something like this:

.. code-block:: bash

    0 0 * * * root /opt/autosnap/autosnap daily
    0 * * * * root /opt/autosnap/autosnap hourly

By default ``autosnap`` will read config from ``/opt/autosnap/autosnap.conf``.
Command line allow one switch ``-c`` to specify alternate configuration file.

One and only one command must be specified in command line. This command must
be the name of interval from configuration file.

During execution, autosnap will create one new snapshot for each included dataset
and delete all oldest snapshots exceeding the allowed snapshots count for given interval.

Additional commands
-------------------

``autosnap`` supports two special commands, ``list-unmanaged-snapshots`` and ``list-managed-snapshots``.

Command ``/opt/autosnap/autosnap list-unmanaged-snapshots`` will list all existing snapshots, which are not managed by ``autosnap``.

Command ``/opt/autosnap/autosnap list-managed-snapshots`` will list all existing snapshots, which are managed by ``autosnap``.

If all snapshots are managed by ``autosnap`` it will be useful to schedule in cron command ``/opt/autosnap/autosnap list-unmanaged-snapshots``
for periodic execution. If abandoned snapshots appears - it will be listed by command ``/opt/autosnap/autosnap list-unmanaged-snapshots``
and report about such abandoned snapshots will be sent to system administrator mail.

If all snapshots are managed by ``autosnap`` and ``autobackup`` cron command ``/opt/autosnap/autosnap list-unmanaged-snapshots | grep -v "@autobackup"``
will be useful.

