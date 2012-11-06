kson
====

Package kson implements encoding and decoding of kson

## overview

* kson, Keep It Simple & Stupid Object Notation
* json-like style, doesn't need quote name with "" and values are separated by \n
* doesn't care data type, doesn't depend on indent(yaml)
* design for some usual scene, such as config file, print readable dump

## kson example 


	Name:	value 		
	Int:	-1024		
	Float:	6.4			
	Bool:	true		
	Date: 	2012-12-21
	String:	abcdefghijklmnopqrstuvwxyz/字符串 #:[]{} 
	Quote:	"[0,1,2,3,4,5,6,7,8,9]"

	Json: 	`
			var father = {
			    "Name": "John",
			    "Age": 32,
			    "Children": [
			        {
			            "Name": "Richard",
			            "Age": 7
			        },
			        {
			            "Name": "Susan",
			            "Age": 4
			        }
			    ]
			};
		`		
	Xml: "	<root>
				<!-- a node -->
				<text>
					I'll be back
				</text>
			</root>

		"
	Empty:	

list example 

	[
		line one
		"[line two]"
		"

		line three

		"
	]


hash example 

	{				
		int:	1024	
		float:	6.4		
		bool:	true	
		string:	string	
		text: 	"
				I'm not a great programmer, 
				I'm a pretty good programmer with great habits.
				"
	} 

compose example

	{	
		Log_Level:	debug
		Listen:		8000

		Roles: [
			{
				Name:	user
				Allow:	[
					/user		
					/order
				]
			} 
			{
				Name:	*				
				Deny: 	[
					/user
					/order
				]
			} 
		]

		Db_Log:	{
			Driver:		mysql			
			Host: 		127.0.0.1
			User:		user
			Password:	password
			Database:	log
		}

		Env:	{
			auth:		http://auth.io
			browser:	ie, chrome, firefox, safari
			key:
		}
	}	


## usage
	
Marshal & Unmarshal example
	
	func example() {
		t := newConfig()

		b, err := kson.Marshal(t)
		if err != nil {
			fmt.Println("kson.Marshal error", err, b)
			return
		}
		fmt.Println(string(b))

		var p Config
		err = kson.Unmarshal(b, &p)
		if err != nil {
			fmt.Println("kson.Unmarshal error", err)
			return
		}

		b, err = json.Marshal(p)
		fmt.Println(string(b))

	}
	
Parse to kson.Node example

	func parse() {
		data := `
		{				
			int:	-1024	
			float:	6.4		
			bool:	true	
			string:	string			
			list:	[
						line 1
						line 2
						line 3
					]
		} 
		`
		node, err := kson.Parse([]byte(data))
		if err != nil {
			fmt.Println("kson.Parse error", err)
			return
		}

		fmt.Println("int", node.ChildInt("int"))
		fmt.Println("float", node.ChildFloat("float"))
		fmt.Println("bool", node.ChildBool("bool"))
		fmt.Println("string", node.ChildString("string"))

		var list []string
		if child, ok := node.Child("list"); ok {
			child.Value(&list)
			fmt.Println("list", list)
		}

		hash := make(map[string]string)
		node.Value(&hash)
		fmt.Println("hash", hash)

	}

    
For more example usage, please see `*_test.go` or `example.go`

## Performance

It should be faster than json

cache reflect.type info maybe improve performance, haven't test it yet

## License

MIT

