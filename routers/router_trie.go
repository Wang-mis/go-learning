package routers

import (
	"fmt"
	"go-koa/middlewares"
	"strings"
)

type RouterTrieNode struct {
	name string
	// Method -> Middleware
	handlers    map[string]middlewares.Middleware
	corsConfigs map[string]*CORSConfig
	children    map[string]*RouterTrieNode
	// 如果当前节点的子节点有以冒号开头的路由，设置该值
	dynamicParam string
}

func newRouterTrieNode(name string) *RouterTrieNode {
	return &RouterTrieNode{
		handlers:    make(map[string]middlewares.Middleware),
		children:    make(map[string]*RouterTrieNode),
		corsConfigs: make(map[string]*CORSConfig),
		name:        name,
	}
}

// newRouterTrie 创建RouterTrie
func newRouterTrie(prefix string) (*RouterTrieNode, *RouterTrieNode) {
	root := newRouterTrieNode("")

	leaf := root
	if prefix[0] == '/' && prefix != "/" {
		leafNode, _ := root.addPath(prefix)
		leaf = leafNode
	}

	return root, leaf
}

// add 向当前节点添加子节点并返回；如果子节点已存在，则直接返回子节点
func (node *RouterTrieNode) add(name string) (*RouterTrieNode, error) {
	if _, exist := node.children[name]; exist {
		return node.children[name], nil
	}

	if name != "" && name[0] == ':' {
		if node.dynamicParam != "" {
			return nil, fmt.Errorf("已存在动态路由参数%s", node.dynamicParam)
		}

		node.dynamicParam = name
	}

	// 创建子节点
	child := newRouterTrieNode(name)
	node.children[name] = child

	return child, nil
}

// addPath 在当前节点添加路径并返回叶节点
func (node *RouterTrieNode) addPath(path string) (*RouterTrieNode, error) {
	nowNode := node

	parts := strings.Split(path[1:], "/")
	for _, part := range parts {
		child, err := nowNode.add(part)
		if err != nil {
			return nil, err
		}
		nowNode = child
	}

	return nowNode, nil
}

// findPath 从当前节点开始寻找路径返回最终叶节点；如果未找到则返回nil
func (node *RouterTrieNode) findPath(path string) (*RouterTrieNode, map[string]string) {
	parts := strings.Split(path[1:], "/")
	params := make(map[string]string)
	nowNode := node
	for _, part := range parts {
		// 先查找静态路由
		child, exists := nowNode.children[part]
		// 静态路由未找到，查找动态路由
		if !exists {
			// 使用动态路由
			if nowNode.dynamicParam != "" {
				params[nowNode.dynamicParam[1:]] = part
				child = nowNode.children[nowNode.dynamicParam]
			} else {
				return nil, params
			}
		}

		nowNode = child
	}

	return nowNode, params
}

// addMethod 当前节点上添加方法处理函数
func (node *RouterTrieNode) addMethod(
	method string,
	handler middlewares.Middleware,
	cors *CORSConfig,
) error {
	if _, exist := node.handlers[method]; exist {
		return fmt.Errorf("方法%s已定义", method)
	}

	node.handlers[method] = handler
	node.corsConfigs[method] = cors
	return nil
}
