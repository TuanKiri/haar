package haar

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The values in this test are taken from the reference C library,
// but they do not match the results produced by this Go package.
//
// Diff:
//
//	--- Expected
//	+++ Actual
//	@@ -2,5 +2,5 @@
//	  Avgl: ([3]float64) (len=3) {
//	-  (float64) 0.25372929123834903,
//	-  (float64) 0.04215254162161761,
//	-  (float64) -0.005735878163064168
//	+  (float64) 0.2529288642406461,
//	+  (float64) 0.04229021787643428,
//	+  (float64) -0.005793594837188721
//	  },
//	@@ -94,3 +94,2 @@
//	    (int16) -1924,
//	-   (int16) -1923,
//	    (int16) -1922,
//	@@ -124,2 +123,3 @@
//	    (int16) 641,
//	+   (int16) 649,
//	    (int16) 768,
func TestSignatureFromBlob(t *testing.T) {
	testCase := struct {
		filename  string
		hash      string
		signature *Signature
	}{
		filename: "testdata/image",
		hash: "iqdb_3fd03d19c7f9fb633fa59504976a8f8dbf777e8111fcba5af4f7fa7cfc00fcfffd00fd77fd" +
			"7ffe7bfe7cfe7dfe7efe7ffe80fef7feffff7cff7dff7eff7fff80fff3fffafffb000100020004000" +
			"70009010202840380038104010403058205890780078107820784ef00f780f9fffa00fc00fcfffd00" +
			"fd77fd7dfe00fe7bfe7efe7ffefeff7aff7dff7fff80fff3fffafffb0001000200030004000700090" +
			"085010301070380038105000581060507800781078207830c05f877f87cf87df87ef87ff880f8f8fa" +
			"77fc7ffc80fd7cfdfffff5fff7fffeffff00030006000d0080008100820083008b018001810182018" +
			"3018401850189028103000301040006800681070007010702",
		signature: &Signature{
			Avgl: [colorChannels]float64{
				0.25372929123834903, 0.04215254162161761, -0.005735878163064168,
			},
			Sig: [colorChannels][haarCoefficients]int16{
				{
					-2825,
					-1412,
					-1024,
					-769,
					-768,
					-649,
					-641,
					-389,
					-388,
					-387,
					-386,
					-385,
					-384,
					-265,
					-257,
					-132,
					-131,
					-130,
					-129,
					-128,
					-13,
					-6,
					-5,
					1,
					2,
					4,
					7,
					9,
					258,
					644,
					896,
					897,
					1025,
					1027,
					1410,
					1417,
					1920,
					1921,
					1922,
					1924,
				},
				{
					-4352,
					-2176,
					-1537,
					-1536,
					-1024,
					-769,
					-768,
					-649,
					-643,
					-512,
					-389,
					-386,
					-385,
					-258,
					-134,
					-131,
					-129,
					-128,
					-13,
					-6,
					-5,
					1,
					2,
					3,
					4,
					7,
					9,
					133,
					259,
					263,
					896,
					897,
					1280,
					1409,
					1541,
					1920,
					1921,
					1922,
					1923,
					3077,
				},
				{
					-1929,
					-1924,
					-1923,
					-1922,
					-1921,
					-1920,
					-1800,
					-1417,
					-897,
					-896,
					-644,
					-513,
					-11,
					-9,
					-2,
					-1,
					3,
					6,
					13,
					128,
					129,
					130,
					131,
					139,
					384,
					385,
					386,
					387,
					388,
					389,
					393,
					641,
					768,
					769,
					1024,
					1664,
					1665,
					1792,
					1793,
					1794,
				},
			},
		},
	}

	image, err := os.ReadFile(testCase.filename)
	require.NoError(t, err)

	blob, err := base64.StdEncoding.DecodeString(string(image))
	require.NoError(t, err)

	signature, err := SignatureFromBlob(blob)
	require.NoError(t, err)

	assert.Equal(t, testCase.signature, signature)
	assert.Equal(t, testCase.hash, signature.String())
}

func TestSignatureFromHash(t *testing.T) {
	testCase := struct {
		hash      string
		signature *Signature
	}{
		hash: "iqdb_3fd03d19c7f9fb633fa59504976a8f8dbf777e8111fcba5af4f7fa7cfc00fcfffd00fd77fd" +
			"7ffe7bfe7cfe7dfe7efe7ffe80fef7feffff7cff7dff7eff7fff80fff3fffafffb000100020004000" +
			"70009010202840380038104010403058205890780078107820784ef00f780f9fffa00fc00fcfffd00" +
			"fd77fd7dfe00fe7bfe7efe7ffefeff7aff7dff7fff80fff3fffafffb0001000200030004000700090" +
			"085010301070380038105000581060507800781078207830c05f877f87cf87df87ef87ff880f8f8fa" +
			"77fc7ffc80fd7cfdfffff5fff7fffeffff00030006000d0080008100820083008b018001810182018" +
			"3018401850189028103000301040006800681070007010702",
		signature: &Signature{
			Avgl: [colorChannels]float64{
				0.25372929123834903, 0.04215254162161761, -0.005735878163064168,
			},
			Sig: [colorChannels][haarCoefficients]int16{
				{
					-2825,
					-1412,
					-1024,
					-769,
					-768,
					-649,
					-641,
					-389,
					-388,
					-387,
					-386,
					-385,
					-384,
					-265,
					-257,
					-132,
					-131,
					-130,
					-129,
					-128,
					-13,
					-6,
					-5,
					1,
					2,
					4,
					7,
					9,
					258,
					644,
					896,
					897,
					1025,
					1027,
					1410,
					1417,
					1920,
					1921,
					1922,
					1924,
				},
				{
					-4352,
					-2176,
					-1537,
					-1536,
					-1024,
					-769,
					-768,
					-649,
					-643,
					-512,
					-389,
					-386,
					-385,
					-258,
					-134,
					-131,
					-129,
					-128,
					-13,
					-6,
					-5,
					1,
					2,
					3,
					4,
					7,
					9,
					133,
					259,
					263,
					896,
					897,
					1280,
					1409,
					1541,
					1920,
					1921,
					1922,
					1923,
					3077,
				},
				{
					-1929,
					-1924,
					-1923,
					-1922,
					-1921,
					-1920,
					-1800,
					-1417,
					-897,
					-896,
					-644,
					-513,
					-11,
					-9,
					-2,
					-1,
					3,
					6,
					13,
					128,
					129,
					130,
					131,
					139,
					384,
					385,
					386,
					387,
					388,
					389,
					393,
					641,
					768,
					769,
					1024,
					1664,
					1665,
					1792,
					1793,
					1794,
				},
			},
		},
	}

	signature, err := SignatureFromHash(testCase.hash)
	require.NoError(t, err)

	assert.Equal(t, testCase.signature, signature)
}
