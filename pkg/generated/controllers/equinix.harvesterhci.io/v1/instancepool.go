/*
Copyright 2022 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/harvester/harvester-equinix-addon/pkg/apis/equinix.harvesterhci.io/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type InstancePoolHandler func(string, *v1.InstancePool) (*v1.InstancePool, error)

type InstancePoolController interface {
	generic.ControllerMeta
	InstancePoolClient

	OnChange(ctx context.Context, name string, sync InstancePoolHandler)
	OnRemove(ctx context.Context, name string, sync InstancePoolHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() InstancePoolCache
}

type InstancePoolClient interface {
	Create(*v1.InstancePool) (*v1.InstancePool, error)
	Update(*v1.InstancePool) (*v1.InstancePool, error)
	UpdateStatus(*v1.InstancePool) (*v1.InstancePool, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.InstancePool, error)
	List(opts metav1.ListOptions) (*v1.InstancePoolList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.InstancePool, err error)
}

type InstancePoolCache interface {
	Get(name string) (*v1.InstancePool, error)
	List(selector labels.Selector) ([]*v1.InstancePool, error)

	AddIndexer(indexName string, indexer InstancePoolIndexer)
	GetByIndex(indexName, key string) ([]*v1.InstancePool, error)
}

type InstancePoolIndexer func(obj *v1.InstancePool) ([]string, error)

type instancePoolController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewInstancePoolController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) InstancePoolController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &instancePoolController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromInstancePoolHandlerToHandler(sync InstancePoolHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.InstancePool
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.InstancePool))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *instancePoolController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.InstancePool))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateInstancePoolDeepCopyOnChange(client InstancePoolClient, obj *v1.InstancePool, handler func(obj *v1.InstancePool) (*v1.InstancePool, error)) (*v1.InstancePool, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *instancePoolController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *instancePoolController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *instancePoolController) OnChange(ctx context.Context, name string, sync InstancePoolHandler) {
	c.AddGenericHandler(ctx, name, FromInstancePoolHandlerToHandler(sync))
}

func (c *instancePoolController) OnRemove(ctx context.Context, name string, sync InstancePoolHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromInstancePoolHandlerToHandler(sync)))
}

func (c *instancePoolController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *instancePoolController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *instancePoolController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *instancePoolController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *instancePoolController) Cache() InstancePoolCache {
	return &instancePoolCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *instancePoolController) Create(obj *v1.InstancePool) (*v1.InstancePool, error) {
	result := &v1.InstancePool{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *instancePoolController) Update(obj *v1.InstancePool) (*v1.InstancePool, error) {
	result := &v1.InstancePool{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *instancePoolController) UpdateStatus(obj *v1.InstancePool) (*v1.InstancePool, error) {
	result := &v1.InstancePool{}
	return result, c.client.UpdateStatus(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *instancePoolController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *instancePoolController) Get(name string, options metav1.GetOptions) (*v1.InstancePool, error) {
	result := &v1.InstancePool{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *instancePoolController) List(opts metav1.ListOptions) (*v1.InstancePoolList, error) {
	result := &v1.InstancePoolList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *instancePoolController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *instancePoolController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v1.InstancePool, error) {
	result := &v1.InstancePool{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type instancePoolCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *instancePoolCache) Get(name string) (*v1.InstancePool, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.InstancePool), nil
}

func (c *instancePoolCache) List(selector labels.Selector) (ret []*v1.InstancePool, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.InstancePool))
	})

	return ret, err
}

func (c *instancePoolCache) AddIndexer(indexName string, indexer InstancePoolIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.InstancePool))
		},
	}))
}

func (c *instancePoolCache) GetByIndex(indexName, key string) (result []*v1.InstancePool, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.InstancePool, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.InstancePool))
	}
	return result, nil
}

type InstancePoolStatusHandler func(obj *v1.InstancePool, status v1.InstancePoolStatus) (v1.InstancePoolStatus, error)

type InstancePoolGeneratingHandler func(obj *v1.InstancePool, status v1.InstancePoolStatus) ([]runtime.Object, v1.InstancePoolStatus, error)

func RegisterInstancePoolStatusHandler(ctx context.Context, controller InstancePoolController, condition condition.Cond, name string, handler InstancePoolStatusHandler) {
	statusHandler := &instancePoolStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromInstancePoolHandlerToHandler(statusHandler.sync))
}

func RegisterInstancePoolGeneratingHandler(ctx context.Context, controller InstancePoolController, apply apply.Apply,
	condition condition.Cond, name string, handler InstancePoolGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &instancePoolGeneratingHandler{
		InstancePoolGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterInstancePoolStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type instancePoolStatusHandler struct {
	client    InstancePoolClient
	condition condition.Cond
	handler   InstancePoolStatusHandler
}

func (a *instancePoolStatusHandler) sync(key string, obj *v1.InstancePool) (*v1.InstancePool, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type instancePoolGeneratingHandler struct {
	InstancePoolGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *instancePoolGeneratingHandler) Remove(key string, obj *v1.InstancePool) (*v1.InstancePool, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.InstancePool{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *instancePoolGeneratingHandler) Handle(obj *v1.InstancePool, status v1.InstancePoolStatus) (v1.InstancePoolStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.InstancePoolGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
