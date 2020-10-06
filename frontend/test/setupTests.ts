import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { cache } from 'swr'

Enzyme.configure({ adapter: new Adapter() });
afterEach(() => {
    cache.clear();
})