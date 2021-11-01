package abookclub

import (
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	testSource := `
<div class="entry">
            <div style="padding-bottom: 14px; margin-bottom: 14px; border-bottom: 1px silver solid;">
                <div class="entry_header_full">
                    <a href="http://abook-club.ru/forum/index.php?showtopic=90170">_Сборник - Хрестоматия 4 класс. Зарубежная литература</a>
                </div>
                <div class="entry_time" style="margin-bottom: 7px;">
                    Naina Kievna, 02.05.2021 19:19
                    —
                    <a class="show_hide link_all" href="#" id="link_90170">свернуть</a>
                </div>

                <div class="entry_content block_all" id="block_90170" style="">
                    
                     <img src="http://gallery.abook-club.ru/var/albums/audiobookcover/Hrestomatiya_4_klass_Zarubezhnaya_literatura.jpg" alt="user posted image" border="0"><br><br><b>Автор:</b> _Сборник<br><b>Название:</b> Зарубежная литература<br><b>Исполнитель:</b> Некрасов Денис; Грачёва Наталья; Гуревич Наталья; Исаев Олег Валерьевич; Старчиков Степан<br><b>Цикл/серия:</b> Хрестоматия 4 класс<br><b>Жанр:</b> Зарубежная классика, Учебная литература<br><b>Издательство:</b> Литрес Паблишинг<br><b>Год издания:</b> 2013<br><b>Качество:</b> mp3, vbr, 64 kbps, 44 kHz, Joint Stereo<br><b>Размер:</b> 86,57 MB<br><b>Длительность:</b> 3:14:12<br><br><b>Описание:</b><br>"Хрестоматия 4 класс"составлен в соответствии со школьной программой. Содержит произведения русского фольклора, русской и зарубежной классики детской литературы. Позволяет в доступной форме ознакомиться с литературными произведениями. Помогает освоить правильное, литературное произношение слов. Увеличивает словарный запас школьника. Помогает учиться читать.<br><br><b>Содержание:</b><br>Ганс Христиан Андерсен<br>Русалочка<br>Оле-Лукойе<br>Свинопас<br><br>Легенды и мифы древней Греции<br>Рождение и воспитание Геракла<br>Немейский лев. Первый подвиг Геракла<br>Лернейская гидра. Второй подвиг Геракла<br>Стимфальские птицы. Третий подвиг Геракла<br>Керинейская лань. Четвертый подвиг Геракла<br>Эриманфский кабан. Пятый подвиг Геракла<br><br>Боги и герои Древнего Рима<br>Эней в подземном мире<br>Как гуси Рим спасли<br><br>Легенды о первых Христианах<br>Сказание о Фёдоре – христианине<br><br><img src="/files/images/cut2.png"><span class="cut">вырезано</span>			<br><a href="https://www.litres.ru/5630246/?lfrom=134352321" title="Купить на ЛитРес" target="_blank">Купить на ЛитРес</a> 

                </div>
                <div class="more" style="text-align: right; font-weight: normal;">
                    <a target="_blank" href="http://abook-club.ru/forum/index.php?showtopic=90170"><img src="/files/images/forum.png" alt="forum">Подробнее на форуме</a>
                </div>
            </div>
        </div>
	`

	reader := strings.NewReader(testSource)

	books, err := parseBody(reader)

	if err != nil {
		t.Error("Error in parse.", err)
	}

	if len(books) != 1 {
		t.Error("No books parsed!")
	}

	book := books[0]

	assertEqual(t, "id", book.Id, "abook-club-90170")
	assertEqual(t, "title", book.Title, "Хрестоматия 4 класс. Зарубежная литература")
	assertEqual(t, "author", book.Author, "_Сборник")
	assertEqual(t, "link", book.Link, "http://abook-club.ru/forum/index.php?showtopic=90170")
	assertEqual(t, "description", book.Description, "&#34;Хрестоматия 4 класс&#34;составлен в соответствии со школьной программой. Содержит произведения русского фольклора, русской и зарубежной классики детской литературы. Позволяет в доступной форме ознакомиться с литературными произведениями. Помогает освоить правильное, литературное произношение слов. Увеличивает словарный запас школьника. Помогает учиться читать.")
	assertEqual(t, "size", book.Size, "86,57 MB")
	assertEqual(t, "quality", book.Quality, "mp3, vbr, 64 kbps, 44 kHz, Joint Stereo")
	assertEqualInt(t, "year", book.Year, 2013)

	var date time.Time
	if !date.Before(book.Date) {
		t.Errorf("Invalid date %v", book.Date)
	}
}

func assertEqual(t *testing.T, field string, s1 string, s2 string) {
	if s1 != s2 {
		t.Errorf("Expected field '%s' to be equal '%s' '%s' .", field, s1, s2)
	}
}

func assertEqualInt(t *testing.T, field string, s1 int, s2 int) {
	if s1 != s2 {
		t.Errorf("Expected field '%s' to be equal '%d' '%d' .", field, s1, s2)
	}
}
