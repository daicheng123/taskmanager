package templates

const OperationTpl = `---
- hosts: "all"
  gather_facts: false
  pre_tasks:
    - name: "transfer scripts"
      copy:
        src: {{ .Src }}
        dest: {{ .Dest }}
        owner: root
        group: root
        mode: 775
  tasks:
    - name: "execute scripts"
      shell: {{ .Command }}
  post_tasks:
    - name: "delete scripts"
      file:
        path: {{ .Dest }}
        state: absent
`
