/*
Copyright 2024 Rancher Labs, Inc.

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
	"sync"
	"time"

	v1 "github.com/rancher/backup-restore-operator/pkg/apis/resources.cattle.io/v1"
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

type BackupHandler func(string, *v1.Backup) (*v1.Backup, error)

type BackupController interface {
	generic.ControllerMeta
	BackupClient

	OnChange(ctx context.Context, name string, sync BackupHandler)
	OnRemove(ctx context.Context, name string, sync BackupHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() BackupCache
}

type BackupClient interface {
	Create(*v1.Backup) (*v1.Backup, error)
	Update(*v1.Backup) (*v1.Backup, error)
	UpdateStatus(*v1.Backup) (*v1.Backup, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Backup, error)
	List(opts metav1.ListOptions) (*v1.BackupList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Backup, err error)
}

type BackupCache interface {
	Get(name string) (*v1.Backup, error)
	List(selector labels.Selector) ([]*v1.Backup, error)

	AddIndexer(indexName string, indexer BackupIndexer)
	GetByIndex(indexName, key string) ([]*v1.Backup, error)
}

type BackupIndexer func(obj *v1.Backup) ([]string, error)

type backupController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewBackupController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) BackupController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &backupController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromBackupHandlerToHandler(sync BackupHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Backup
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Backup))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *backupController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Backup))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateBackupDeepCopyOnChange(client BackupClient, obj *v1.Backup, handler func(obj *v1.Backup) (*v1.Backup, error)) (*v1.Backup, error) {
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

func (c *backupController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *backupController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *backupController) OnChange(ctx context.Context, name string, sync BackupHandler) {
	c.AddGenericHandler(ctx, name, FromBackupHandlerToHandler(sync))
}

func (c *backupController) OnRemove(ctx context.Context, name string, sync BackupHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromBackupHandlerToHandler(sync)))
}

func (c *backupController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *backupController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *backupController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *backupController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *backupController) Cache() BackupCache {
	return &backupCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *backupController) Create(obj *v1.Backup) (*v1.Backup, error) {
	result := &v1.Backup{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *backupController) Update(obj *v1.Backup) (*v1.Backup, error) {
	result := &v1.Backup{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *backupController) UpdateStatus(obj *v1.Backup) (*v1.Backup, error) {
	result := &v1.Backup{}
	return result, c.client.UpdateStatus(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *backupController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *backupController) Get(name string, options metav1.GetOptions) (*v1.Backup, error) {
	result := &v1.Backup{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *backupController) List(opts metav1.ListOptions) (*v1.BackupList, error) {
	result := &v1.BackupList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *backupController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *backupController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v1.Backup, error) {
	result := &v1.Backup{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type backupCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *backupCache) Get(name string) (*v1.Backup, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.Backup), nil
}

func (c *backupCache) List(selector labels.Selector) (ret []*v1.Backup, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Backup))
	})

	return ret, err
}

func (c *backupCache) AddIndexer(indexName string, indexer BackupIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Backup))
		},
	}))
}

func (c *backupCache) GetByIndex(indexName, key string) (result []*v1.Backup, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Backup, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Backup))
	}
	return result, nil
}

// BackupStatusHandler is executed for every added or modified Backup. Should return the new status to be updated
type BackupStatusHandler func(obj *v1.Backup, status v1.BackupStatus) (v1.BackupStatus, error)

// BackupGeneratingHandler is the top-level handler that is executed for every Backup event. It extends BackupStatusHandler by a returning a slice of child objects to be passed to apply.Apply
type BackupGeneratingHandler func(obj *v1.Backup, status v1.BackupStatus) ([]runtime.Object, v1.BackupStatus, error)

// RegisterBackupStatusHandler configures a BackupController to execute a BackupStatusHandler for every events observed.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterBackupStatusHandler(ctx context.Context, controller BackupController, condition condition.Cond, name string, handler BackupStatusHandler) {
	statusHandler := &backupStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromBackupHandlerToHandler(statusHandler.sync))
}

// RegisterBackupGeneratingHandler configures a BackupController to execute a BackupGeneratingHandler for every events observed, passing the returned objects to the provided apply.Apply.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterBackupGeneratingHandler(ctx context.Context, controller BackupController, apply apply.Apply,
	condition condition.Cond, name string, handler BackupGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &backupGeneratingHandler{
		BackupGeneratingHandler: handler,
		apply:                   apply,
		name:                    name,
		gvk:                     controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterBackupStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type backupStatusHandler struct {
	client    BackupClient
	condition condition.Cond
	handler   BackupStatusHandler
}

// sync is executed on every resource addition or modification. Executes the configured handlers and sends the updated status to the Kubernetes API
func (a *backupStatusHandler) sync(key string, obj *v1.Backup) (*v1.Backup, error) {
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

type backupGeneratingHandler struct {
	BackupGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
	seen  sync.Map
}

// Remove handles the observed deletion of a resource, cascade deleting every associated resource previously applied
func (a *backupGeneratingHandler) Remove(key string, obj *v1.Backup) (*v1.Backup, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.Backup{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	if a.opts.UniqueApplyForResourceVersion {
		a.seen.Delete(key)
	}

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

// Handle executes the configured BackupGeneratingHandler and pass the resulting objects to apply.Apply, finally returning the new status of the resource
func (a *backupGeneratingHandler) Handle(obj *v1.Backup, status v1.BackupStatus) (v1.BackupStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.BackupGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}
	if !a.isNewResourceVersion(obj) {
		return newStatus, nil
	}

	err = generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
	if err != nil {
		return newStatus, err
	}
	a.storeResourceVersion(obj)
	return newStatus, nil
}

// isNewResourceVersion detects if a specific resource version was already successfully processed.
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *backupGeneratingHandler) isNewResourceVersion(obj *v1.Backup) bool {
	if !a.opts.UniqueApplyForResourceVersion {
		return true
	}

	// Apply once per resource version
	key := obj.Namespace + "/" + obj.Name
	previous, ok := a.seen.Load(key)
	return !ok || previous != obj.ResourceVersion
}

// storeResourceVersion keeps track of the latest resource version of an object for which Apply was executed
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *backupGeneratingHandler) storeResourceVersion(obj *v1.Backup) {
	if !a.opts.UniqueApplyForResourceVersion {
		return
	}

	key := obj.Namespace + "/" + obj.Name
	a.seen.Store(key, obj.ResourceVersion)
}
