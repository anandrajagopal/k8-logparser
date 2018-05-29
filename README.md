# k8-logparser

This is a small utility that runs as a side car to parse logs from a main container. Tools like fluentd parse each line separately and this causes issues when an application runs inside a K8 pod and writes logs to stdout. For example, if java verbose gc and application logs (and potentially other kinds of logs) go to stdout, it is difficult to know how many lines of logs need to be grouped together before logs can be shipped to ELK. This tool/toy program/proof of concept attempts to solve it by having the application write the logs to a named pipe. Each type of log would be written to a named pipe and this side car container (it's configurable) reads from those pipes concurrently using the specified regex. The parser starts by reading from the named pipe and starts accumulating lines until the specified regex matches. Once the regex matches, the parser adds '{{' to the beginning and '}}' to the end of the accumulated lines and pushes it to stdout. At this point, tools like fluentd can be configured to read from the stdout of the sidecar. This method is language independent. 

# How to use

Look at the logparser.yaml for an example. The example uses a springboot app that writes to a named pipe called ouptput.log. This pipe is created by the init container and is mounted at /src/greeting directory on the main application container and /etc/log in the side car container. The GC logs are written to another named pipe called gclog. This pipe is mounted at /etc/gc in the main application container and at the same location in the side car container. These are all configurable. 

The log parser side car container knows about the named pipes and its location using a configmap that is mounted on to it. 


