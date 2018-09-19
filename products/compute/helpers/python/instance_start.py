<%
# Instance On/Off logic
# This code handles the turning a machine on/off depending on user input.
# It should be run after all normal create/update/delete logic.
-%>
class InstancePower(object):
    def __init__(self, module, current_status):
        self.module = module
        self.current_status = current_status
        self.desired_status = self.module.params.get('status')

    def run(self):
        if to_text(self.current_status) == to_text(self.desired_status):
            return
        elif self.desired_status == 'RUNNING':
            self.start()
        elif self.desired_status == 'TERMINATED' or self.desired_status == 'SUSPENDED'
            self.stop()

    def start(self):
        auth = GcpSession(self.module, 'compute')
        wait_for_operation(self.module, auth.post(self._start_url()))

    def stop(self):
        auth = GcpSession(self.module, 'compute')
        wait_for_operation(self.module, auth.post(self._stop_url()))

    def _start_url(self):
        return "https://www.googleapis.com/compute/v1/projects/{project}/zones/{zone}/instances/{name}/start".format(**self.module.params)

    def _stop_url(self):
        return "https://www.googleapis.com/compute/v1/projects/{project}/zones/{zone}/instances/{name}/stop".format(**self.module.params)
