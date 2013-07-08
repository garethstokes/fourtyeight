fourtyeight ( aka, the drop )
=============================

+ setting up a dev environment
------------------------------

go install vagrant ( www.vagrantup.com ) then

```bash
  # this might take a while, 30 mins or so
  vagrant up

  # remote into the machine
  vagrant ssh

  # cd into the src directory
  # i've made a alias 'sw' which stands for 'start work' to do this
  cd /golang/src/github.com/garethstokes/fourtyeight

  # we need to create databases and seed them
  ./prepare_application.sh

  # run the web app
  ./run_website.sh
```

open a new terminal and type in...
curl localhost:8080/

"let thy object decend as if it were calescent"

vagrant will forward port 80 on the virtual machine to your
local host environment. 

nginx is installed on the vm, it will forward the port 8080
and host it on port 80
