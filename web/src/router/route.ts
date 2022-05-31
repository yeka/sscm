"use strict"

export { Tree }

class Tree {
    childs: {[key: string]: Tree} = {}
    value: any
    
    constructor(routes: {[key: string]: any}) {
        for (const path in routes) {
            this.addPath(path, routes[path])
        }
    }

    add(path: string[], val: any): void {
        if (path.length == 0) {
        this.value = val
        return
        }
        if (this.childs[path[0]] == undefined) {
        this.childs[path[0]] = new Tree({});
        }
        this.childs[path[0]].add(path.slice(1, path.length), val)
    }
    
    find(path: string[]): any {
        if (path.length == 0) {
        return {value: this.value, params: {}}
        }
        var c = this.childs[path[0]];
        var params: {[key: string]: string} = {}
        if (c == undefined) {
        for (const k in this.childs) {
            if (k[0] == ":") {
            c = this.childs[k]
            params[k.substring(1)] = path[0]
            break
            }
        }
        }
        if (c == undefined) {
        return undefined
        }
        const res = c.find(path.slice(1, path.length))
        if (res == undefined) {
        return undefined
        }
        params = {...params, ...res.params}
        return {value: res.value, params: params}
    }
    
    addPath(path: string, val: any): void {
        if (path == "" || path == "/") {
            this.value = val
            return
        }
        if (path[0] == "/") {
            path = path.substring(1)
        }
        const parts = path.split("/")
        this.add([parts.length + "", ...parts], val)
    }
    
    findPath(path: string): any {
        if (path == "" || path == "/") {
            return {value: this.value, params: {}}
        }
        if (path[0] == "/") {
            path = path.substring(1)
        }
        const parts = path.split("/")
        return this.find([parts.length + "", ...parts])
    }
}
