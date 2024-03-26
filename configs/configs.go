package configs

import (
	"embed"
	"os"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	myi18n "github.com/raulaguila/go-template/internal/pkg/i18n"
	"github.com/raulaguila/go-template/pkg/helpers"
	"golang.org/x/text/language"
)

//go:embed i18n/*
var viewsfs embed.FS

//go:embed version.txt
var version string

func init() {
	err := godotenv.Load(path.Join("configs", ".env"))
	helpers.PanicIfErr(err)

	os.Setenv("SYS_VERSION", version)

	time.Local, err = time.LoadLocation(os.Getenv("TZ"))
	helpers.PanicIfErr(err)
	helpers.PanicIfErr(loadMessages())
}

func loadMessages() error {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, lang := range strings.Split(os.Getenv("SYS_LANGUAGES"), ",") {
		if _, err := bundle.LoadMessageFileFS(viewsfs, path.Join("i18n", "active."+lang+".toml")); err != nil {
			return err
		}

		myi18n.I18nTranslations[lang] = myi18n.NewTranslation(i18n.NewLocalizer(bundle, lang))
	}

	return nil
}
