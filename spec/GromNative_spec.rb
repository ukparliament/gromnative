RSpec.describe GromNative do
  it 'has a version number' do
    expect(GromNative::VERSION).not_to be nil
  end

  describe '.get' do
    context 'with a valid url' do
      context 'and a 200 status code' do
        it 'makes the expected request' do
          subject.get('http://localhost:3333/full.nt')
        end

        pending 'add tests'
      end

      context 'and a 5xx status code' do
        it { subject.get('http://localhost:3333/500') }
        pending 'add tests'
      end

      context 'and a 3xx status code' do
        it { subject.get('http://localhost:3333/302') }
        pending 'add tests'
      end

      context 'and a 4xx status code' do
        it { subject.get('http://localhost:3333/404') }
        pending 'add tests'
      end
    end

    context 'with an invalid url' do
      it 'returns the expected object' do
        expect(JSON.parse(subject.get('foo://a_broken.url'))).to eq({'statementsBySubject' => nil,'edgesBySubject' => nil,'statusCode' => 0,'uri' => '','error' => "Error getting data: Get foo://a_broken.url: unsupported protocol scheme \"foo\"\n"})
      end
    end
  end
end
