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
by name or by pattern. Pattern is rsync-like, ``?`` means any symbol,
``*`` means any symbol except ``/`` symbol, ``**`` means any symbol.

First match win, and if it was directive ``exclude`` - dataset will be excluded,
if it was directive ``include`` - dataset will be included.

Schedule autosnap
-----------------
  - ``vim /etc/cron.d/autosnap``
  - write to config something like this:

.. code-block:: bash

    0 0 * * * root /opt/autosnap/autosnap daily
    0 * * * * root /opt/autosnap/autosnap hourly

By default ``autosnap`` will read config from ``/opt/autosnap/autosnap.conf``.
Command line allow one switch ``-c`` to specify alternate configuration file.

One and only one command must be specified in command line. This command must
be the name of interval from configuration file.

During execution, autosnap will create one new snapshot for each included dataset
and will delete all oldest snapshots exceeding the allowed snapshots count for given interval.

