# Macros in LMV-Suite - Complete Guide

Macros are powerful shortcuts that let you save and reuse commands quickly.  
You can create your own macros and override built-in ones.

## 1. Defining Macros

Basic syntax (both forms are accepted):

```sh
#def      name |param1,param2,...|          -> command template
#define   name |param1:must,param2|        -> command template
```

### Parameter rules:
- Parameters are listed inside `|` ... `|`
- Separate parameters with commas
- Add `:must` after a parameter name to make it **required**
- Use `$param_name` in the command to insert the value

### Definition examples:

```sh
#def greet        |name|                  -> #echo Hello $name!
#def scan         |target:must,port|      -> nmap -p $port -T4 $target
#def set_target   |ip:must|               -> #set current_target $ip
#def check_port   |host,port=80|          -> nc -zv $host $port
#def backup       |path:must,verbose|     -> tar -czf backup.tar.gz $path
```

## 2. All Ways to Call / Use Macros

Macros support multiple calling styles — choose whichever feels most natural:

| Style                                 | Example                                      | Description                              |
|---------------------------------------|----------------------------------------------|------------------------------------------|
| Simple positional                     | `#greet hamza`                               | Most common - first value = first param  |
| Quoted positional                     | `#greet "Mr. Smith"`                         | Use quotes for values with spaces        |
| With parentheses                      | `#greet("alice")`                            | Clean style, good for single argument    |
| Named parameter                       | `#scan target=192.168.1.100`                 | Explicit - good when multiple params     |
| Named with parentheses                | `#scan(target="10.10.10.50")`                | Very explicit style                      |
| Mixed style                           | `#backup path=/home verbose`                 | Positional + optional named              |
| Multiple arguments                    | `#connect host=server user=root pass=secret` | Real-world complex call                  |

All these are valid and will work the same:

```sh
#greet hamza
#greet "hamza khan"
#greet("alice")
#greet name="david"
#greet(name="sara")
#scan target=192.168.1.100 port=22
#scan 192.168.1.100 80
```

## 3. Built-in Macros (always available)

These are pre-defined and ready to use. You can override them with your own `#def` if needed.

```sh
#echo               text here                → prints text
#if cond -> command → simple conditional (cond: true/1/yes/on)
#pwd                                         → current working directory
#whoami                                      → current username
#date                                        → current date & time
#clear    or  #cls                           → clear screen
#value    $VAR                               → show value of environment variable
```

## 4. Real-World Practical Examples

```sh
#def target       |ip:must|              -> #set current_target $ip
#def quick        |port|                 -> nmap -sV -T4 -p $port $current_target
#def full         |port=1-65535|         -> nmap -A -T4 -p $port $current_target
#def users                             -> cat /etc/passwd | cut -d: -f1 | sort

#target 10.10.10.50
#quick 80
#full
#users
```

## 5. Important Notes

- **Required parameters** (`:must`) → error if missing
- Missing required param example error: