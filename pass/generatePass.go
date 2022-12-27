package pass

import (
	"github.com/gin-gonic/gin"
	passkit "github.com/ksaduakasov/kalbekpasskit"
	"net/http"
	"net/url"
	"os"
)

func GeneratePass(c *gin.Context, t string, p string, o string, s string, k string, l string, v string, icon string, logo string, strip string) (string, error) {
	storeCard := passkit.NewStoreCard()

	primaryField := passkit.Field{
		Key:   "Primary key",
		Label: "Primary label",
		Value: "Primary value",
	}

	storeCard.AddPrimaryFields(primaryField)

	secondaryField := passkit.Field{
		Key:   k,
		Label: l,
		Value: v,
	}

	storeCard.AddSecondaryFields(secondaryField)

	pass := passkit.Pass{
		FormatVersion:       1,
		TeamIdentifier:      t,
		PassTypeIdentifier:  p,
		OrganizationName:    o,
		SerialNumber:        s,
		Description:         "Store Card for asd",
		StoreCard:           storeCard,
		LogoText:            "Cleverest Technologies",
		WebServiceURL:       "http://localhost:3000/",
		AuthenticationToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIwNTMxMDYsInN1YiI6MTB9.-bZVaodJigPANxHaEcEWykyl-MWDTjDnz8zxCuVk8tk",
		Barcodes: []passkit.Barcodes{
			{
				Format:          passkit.BarcodeFormatQR,
				Message:         "22618981-46e9-4322-8ff4-84b49921073a",
				MessageEncoding: "utf-8",
			},
		},
	}

	template := passkit.NewInMemoryPassTemplate()
	iconURL, err := url.Parse(icon)
	if err != nil {
		panic(err)
	}
	logoURL, err := url.Parse(logo)
	if err != nil {
		panic(err)
	}
	stripURL, err := url.Parse(strip)
	if err != nil {
		panic(err)
	}

	template.AddFileFromURL(passkit.BundleIcon, url.URL{
		Scheme: iconURL.Scheme,
		Host:   iconURL.Host,
		Path:   iconURL.Path,
	})
	template.AddFileFromURL(passkit.BundleLogo, url.URL{
		Scheme: logoURL.Scheme,
		Host:   logoURL.Host,
		Path:   logoURL.Path,
	})
	template.AddFileFromURL(passkit.BundleStrip, url.URL{
		Scheme: stripURL.Scheme,
		Host:   stripURL.Host,
		Path:   stripURL.Path,
	})

	signer := passkit.NewMemoryBasedSigner()

	signInfo, err := passkit.LoadSigningInformationFromFiles("/Users/ksaduakassov/Documents/Cleverest-cer-pass.p12", "Cleverest2022", "/Users/ksaduakassov/Documents/AppleWWDRCAG6.cer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return err.Error(), err
	}

	zippedPass, err := signer.CreateSignedAndZippedPassArchive(&pass, template, signInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return err.Error(), err
	}

	err = os.WriteFile("./pass.pkpass", zippedPass, 0644) //needs to be stored somewhere
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return err.Error(), err
	}
	return "success", nil
}
