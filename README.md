[Agni](https://en.wikipedia.org/wiki/Agni) is a bridge between the Prometheus
monitoring system and the [RightLink](http://docs.rightscale.com/rl10/) system
management agent. It can be invoked as a collectd plugin to emit monitoring
series for long-term storage, alerting and other functions in RightScale's
[dashboard](https://login.rightscale.com).

## Configuration

There are two configuration files, both located in `/etc/agni`. Querier
configuration is optional if you want to use the defaults.

Example `plugins.yaml`:

```yaml
magic:
  gauge-magic_smoke_level: rate(go_memstats_alloc_bytes[20s])
  counter-magic_bunnies: go_memstats_frees_total
```

Example `querier.yaml` assuming that Prometheus is not running at the default
location (`http://localhost:9090`):

```yaml
  url: http://127.0.0.1:32770
```

## Integration with RightLink

http://docs.rightscale.com/rl10/reference/10.5.3/rl10_monitoring.html#custom-monitoring-plugins-with-built-in-monitoring

Copy the plugin to a remote machine:

```bash
scp agni 34.203.244.119:agni
```

Make sure `/etc/agni` contains `plugins.yaml` and `querier.yaml`. Tell RightLink
where to find the plugin and make sure all monitoring features are enabled:

```bash
ssh -L 9090:localhost:32771 34.203.244.119

sudo rsc rl10 create /rll/tss/exec/agni executable=/home/rightscale12853/agni

sudo rsc rl10 show /rll/tss/control # NB: copy-paste & set TSS_ID by hand!
sudo rsc rl10 put_control /rll/tss/control enable_monitoring=all tss_id=$TSS_ID
```

Removing:

```bash
sudo rsc rl10 destroy /rll/tss/exec/agni
```
