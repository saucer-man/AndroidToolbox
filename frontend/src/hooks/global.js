import {getCurrentInstance} from 'vue'

export default function useGetGlobalProperties() {
	const {emit, appContext: {app: {config: {globalProperties}}}} = getCurrentInstance()
	
	return {...globalProperties}
}