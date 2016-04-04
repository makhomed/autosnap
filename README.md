# autosnap
ZFS snapshots automation tool
## Usage
* [install Go](https://golang.org/doc/install)
  - [Download Go](https://golang.org/dl/)
  - tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
  - save ```export PATH=$PATH:/usr/local/go/bin``` to ```/etc/profile.d/go.sh```
* Install autosnap
  - ```mkdir /opt/autosnap```
  - ```mkdir /opt/autosnap/bin```
  - ```mkdir /opt/autosnap/conf```
  - ```cd /opt/autosnap```
  - ```git clone https://github.com/makhomed/autosnap.git src```
  - ```cd src```
  - ```export GOPATH=/opt/autosnap```
  - ```export GOBIN=/opt/autosnap/bin```
  - ```/usr/local/go/bin/go install autosnap.go```
* Configure autosnap
  - ```cp /opt/autosnap/src/example/autosnap.conf /opt/autosnap/conf/autosnap.conf```
  - ```vim /opt/autosnap/conf/autosnap.conf```
* Schedule autosnap
  - ```cp /opt/autosnap/src/example/autosnap.cron /etc/cron.d/autosnap```
  - ```vim /etc/cron.d/autosnap```
* Advanced usage
  - ```/opt/autosnap/bin/autosnap clean``` will only delete old ```autosnap``` snapshots.
