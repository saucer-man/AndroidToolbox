
import {getCurrentInstance} from 'vue'

export default function useGlobalProperties() {
	const {emit, appContext:{app:{config: {globalProperties}}}}= getCurrentInstance()
	return {...globalProperties}
}


