operator: {
	"- hosts": "all"
		 gather_facts: "false"
		 pre_tasks: [
			{
				name: params.transferName
				copy: {
					src: params.src
					dest: params.dest
					owner: params.owner
					group: params.group
					mode: params.mode
						}
				},
		]
		tasks: [
			{
				name: params.executeName
				shell: params.command
					}
			]
		post_tasks: [
	  	{
	  		name: params.deleteFileName
	  		file:
	  	  path: params.dest
	  	  state: "absent"
	  	}
		]
	}

params: {
    executeName: string | *"execute script"
    transferName: string | *"transfer script"
    deleteFileName: string | *"delete script"
    owner: string |*"root"
    group: string |*"root"
    mode: int | *775
    src: string
    dest: string
  	command: string
}