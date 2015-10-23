# Goscli
The *Goscli* package is a cross-platform CLI tool for ScaleIO.  The tool is built on top of the *Goscaleio* package which is a reusable implementation of API bindings for ScaleIO written in Go.


- [Current State](#state)
- [Usage](#usage)
- [Licensing](#licensing)
- [Support](#support)

## <a id="state">Current State</a>

There are plenty of great features of the CLI, but it is still early stages.  The capabilities are listed below.

For now, reference the [Goair project](https://github.com/emccode/goair) for basics relating to configuration files, environment variables, and Docker containers as they are built using the same Go CLI framework and packages.


## <a id="usage">Usage</a>

    goscli login --username=admin --password=Scaleio123 --endpoint=https://ip_or_dns_of_gw/api --insecure=true
    goscli instance get
    goscli system use --systemid=38a6603e69c6b8b1
    goscli system get
    goscli statistics get
    goscli user get
    goscli protectiondomain get
    goscli scsiinitiaor get
    goscli sdc get
    goscli sdc get --sdcguid=60424D25-BA83-4324-8E6D-3CED74FB2A30
    goscli sdc local
    goscli sdc local volume
    goscli sdc local statistics
    goscli sdc get --sdcid=988a23eb00000002 volume
    goscli protectiondomain use --protectiondomainid=dbe9a4b700000000
    goscli storagepool get
    goscli storagepool use --storagepoolid=aab7ee0800000001
    goscli volume get
    goscli volume local
    goscli volume get snapshot --volumename=testing3
    goscli volume get vtree --volumename=testing3
    goscli volume get --ancestorvolumeid=d2a3950700000007
    goscli volume create --volumename=testing4 --volumesizeinkb=4096
    goscli volume map --sdcid=988a23eb00000002 --volumeid=d2a3950700000007
    goscli volume map local --volumename=testing1
    goscli volume unmap --sdcid=988a23eb00000002 --volumeid=d2a3950700000007
    goscli volume unmap local --volumename=testing1
    goscli volume remove --volumeid=d2a3950700000007
    goscli volume remove --ancestorvolumeid=d2a3950700000028
    goscli volume remove-snapshot --volumeid=d2a3950700000028



<a id="licensing">Licensing</a>
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

<a id="support">Support</a>
-------
Please file bugs and issues on the Github issues page for this project. This is to help keep track and document everything related to this repo. For general discussions and further support you can join the [EMC {code} Community slack channel](http://community.emccode.com/). Lastly, for questions asked on [Stackoverflow.com](https://stackoverflow.com) please tag them with **EMC**. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.

