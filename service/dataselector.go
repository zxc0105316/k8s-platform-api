package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

//  数据过滤

// 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDatalist []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell接口，用于各种资源的List的类型转换，转换后可以使用dataSelector的排序、过滤、分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// 定义过滤和分页的结构体
type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

// 通过名字进行过滤
type FilterQuery struct {
	Name string
}

// 分页两个属性  limit  和  page
type PaginateQuery struct {
	Limit int
	Page  int
}

// 实现自定义结构的排序，需要重写len,swap,less方法
// len方法用户获取数组的长度

func (d *dataSelector) Len() int {
	return len(d.GenericDatalist)
}

// swap方法用于数据比较大小之后的位置变更
func (d *dataSelector) Swap(i, j int) {
	// python的对换方式
	d.GenericDatalist[i], d.GenericDatalist[j] = d.GenericDatalist[j], d.GenericDatalist[i]

}

// Less方法用于比较大小
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDatalist[i].GetCreation()
	b := d.GenericDatalist[j].GetCreation()

	return b.Before(a)
}

// 排序    上面重写的len , swap , less 都是为下面这个sort功能服务的
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 根据Name过滤数据
// Filter方法用于过滤数据，比较数据的Name属性，若包含，则返回
func (d *dataSelector) Filter() *dataSelector {
	//	 判断入参是否为空，若为空，则返回所有数据
	if d.DataSelect.Filter.Name == "" {
		return d
	} else {
		//	 若不为空，则按照入参Name进行过滤
		// 	 声明一个新的数组，若Name包含，则把数据放进数组，返回出去
		filterArray := []DataCell{}

		for _, value := range d.GenericDatalist {
			// 定义是否匹配的标签变量，默认是匹配
			matchs := true
			objName := value.GetName()
			if !strings.Contains(objName, d.DataSelect.Filter.Name) {
				matchs = false
				// continue 进入下一次循环
				continue
			}
			if matchs {
				filterArray = append(filterArray, value)
			}
		}
		// 将过滤好的数据写入
		d.GenericDatalist = filterArray
	}
	return d
}

// 分页逻辑  根据limit 和 page的传参，取一定范围内的数据

func (d *dataSelector) Paginate() *dataSelector {
	// 根据limit和page的入参，定义快捷变量

	limit := d.DataSelect.Paginate.Limit
	page := d.DataSelect.Paginate.Page
	// 检验参数的合法性
	if limit <= 0 || page <= 0 {
		return d
	}

	//定义取数范围需要的 startIndes 和 endIndex

	// 1 , 10  ,11 20  21 end
	startIndex := limit * (page - 1)
	endIndex := limit*page - 1

	// 处理endIndex，当endIndex大于数组长度,让endIndex等于数组长度
	if endIndex > len(d.GenericDatalist) {
		endIndex = len(d.GenericDatalist)
	}
	// 分段号的数据
	d.GenericDatalist = d.GenericDatalist[startIndex:endIndex]
	return d
}

// 定义podCell ,重写GetCreation 和 GetName方法，可以进行数据转换
type podCell corev1.Pod

// covev1.pod > podCell > DataCell
// appsv1.deployment > deployCell > DataCell

// 重写DataCell 接口的两个方法

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

// deployment

type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

// service

type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}

// pv
type pvtCell corev1.PersistentVolume

func (pv pvtCell) GetCreation() time.Time {
	return pv.CreationTimestamp.Time
}

func (pv pvtCell) GetName() string {
	return pv.Name
}

// namespace
type NamespaceCell corev1.Namespace

func (n NamespaceCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n NamespaceCell) GetName() string {
	return n.Name
}

// statefulset
type statefulsetCell appsv1.StatefulSet

func (s statefulsetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulsetCell) GetName() string {
	return s.Name
}

// daemonset

type daemonsetCell appsv1.DaemonSet

func (d daemonsetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonsetCell) GetName() string {
	return d.Name
}

// ingress

type ingressCell nwv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}

// PersistentVolumes

type persistentVolumeCell corev1.PersistentVolume

func (p persistentVolumeCell) GetCreation() time.Time {

	return p.CreationTimestamp.Time

}

func (p persistentVolumeCell) GetName() string {

	return p.Name
}

// secret

type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {

	return s.CreationTimestamp.Time

}

func (s secretCell) GetName() string {

	return s.Name
}
